package setup

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func Encoder() {
	go func() {
		// update old states of histories on startup
		var histories []m.History
		if err := app.DB.
			Where("status != ?", "failed").
			Where("status != ?", "finished").
			Find(&histories).Error; err != nil {
			log.Error("Failed to list old histories: ", err)
			return
		}
		for _, history := range histories {
			historyPtr := &history
			if err := historyPtr.Failed(app.DB, "Failed because serever shutdown (probably)"); err != nil {
				log.Error("Failed to upate history: ", err)
			}
		}
	}()
	go func() {
		for {
			if !app.Setting.EnableEncoding {
				time.Sleep(time.Second * 5)
				continue
			}
			if len(app.FilesToEncode) > 0 {
				fileToEncode := app.FilesToEncode[0]
				if app.Setting.EncodingMaxRetry > 0 {
					hash, err := helper.HashFile(fileToEncode)
					log.Debug("Failed to hash file to encode", err)
					if err != nil {
						app.FilesToEncode = app.FilesToEncode[1:]
						continue
					}
					var tries int64
					if err := app.DB.
						Model(&m.History{}).
						Where(&m.History{Hash: hash}).
						Count(&tries).Error; err != nil {
						log.Error("Failed to count encoding tries: ", err)
						app.FilesToEncode = app.FilesToEncode[1:]
						continue
					}
					if tries >= int64(app.Setting.EncodingMaxRetry) {
						log.Debug("Reached max retries of file ", fileToEncode)
						app.FilesToEncode = app.FilesToEncode[1:]
						continue
					}
				}
				encodeFile(fileToEncode)
				app.FilesToEncode = app.FilesToEncode[1:]
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func encodeFile(file string) {
	app.CurrentFileToEncode = file
	defer func() {
		app.CurrentFileToEncode = ""
	}()

	// TODO: legacy logic
	if strings.Contains(file, "[encoded]") {
		log.Infof("Skipping already encoded file %s\n", file)
		return
	}

	history := &m.History{}
	if err := history.Create(app.DB, file); err != nil {
		log.Errorf("Failed to create history %v\n", err)
		return
	}

	log.Infof("Encoding file %s\n", file)
	fi, err := os.Stat(file)
	if err != nil {
		log.Errorf("Failed to read filesize %v\n", err)
		if err := history.Failed(app.DB, err.Error()); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}
	oldSize := fi.Size()
	if err := history.SetOldSize(app.DB, uint64(oldSize)); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

	hash, err := helper.HashFile(file)
	if err != nil {
		if err := history.Failed(app.DB, err.Error()); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}
	if err := history.SetHash(app.DB, hash); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

	output := strings.TrimSuffix(file, ".mkv")
	output = fmt.Sprintf("%s[encoded]%s", output, ".mkv")
	tmpOutput := "tmp.mkv"

	defer os.Remove(tmpOutput)

	if err := history.SetNewPath(app.DB, output); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

	if err := history.Encoding(app.DB); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

	// probe file so we can show encoding progress
	// ffprobe context
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	// probe file
	data, err := ffprobe.ProbeURL(ctx, file)
	if err != nil {
		if err := history.Failed(app.DB, err.Error()); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}
	dataStreams := data.StreamType(ffprobe.StreamAny)
	videoDuration := data.Format.Duration().Seconds()
	hasVideoStream := false

	// loop over streams in file
	for _, streamInfo := range dataStreams {
		if streamInfo.CodecType == "video" {
			hasVideoStream = true
		}
	}

	if !hasVideoStream {
		if err := history.Failed(app.DB, "No video stream detected"); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}

	// https://www.tauceti.blog/posts/linux-ffmpeg-amd-5700xt-hardware-video-encoding-hevc-h265-vaapi/
	// https://trac.ffmpeg.org/ticket/3730
	// https://x265.readthedocs.io/en/latest/cli.html#performance-options
	var ffmpegCommand string
	if app.Setting.EnableHevcEncoding {
		h265Pools := "*"
		if app.Setting.EncodingThreads > 0 {
			h265Pools = fmt.Sprint(app.Setting.EncodingThreads)
		}
		ffmpegCommand =
			"ffmpeg " +
				// "-analyzeduration 30000000 -probesize 8000000000 " +
				fmt.Sprintf(`-i "%s" `, file) + // input file
				// "-max_muxing_queue_size 9999 " +
				fmt.Sprintf("-threads %d ", app.Setting.EncodingThreads) +
				"-c:a copy " +
				"-c:s copy " +
				"-c:v libx265 " + // setting video codec libx265 | libaom-av1
				"-map 0:v:0 " +
				"-map 0:a? " +
				"-map 0:s? " +
				// "-pix_fmt yuv420p " + // YUV 4:2:0
				"-profile:v main " + // force 8 bit
				fmt.Sprintf("-crf %d ", app.Setting.EncodingCrf) + // setting quality
				fmt.Sprintf("-x265-params crf=%d:pools=%s -strict experimental ", app.Setting.EncodingCrf, h265Pools) + // setting libx265 params
				fmt.Sprintf("-filter:v scale=%d:-2 ", app.Setting.EncodingResolution) + // setting resolution
				fmt.Sprintf(`"%s" `, tmpOutput) +
				fmt.Sprintf("-progress unix://%s -y", tempSock(
					videoDuration,
					fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
					history,
				)) // progress tracking
	} else {
		ffmpegCommand =
			"ffmpeg " +
				fmt.Sprintf(`-i "%s" `, file) + // input file
				fmt.Sprintf("-threads %d ", app.Setting.EncodingThreads) +
				"-c:a copy " +
				"-c:s copy " +
				"-c:v libx264 " + // setting video codec libx264 | libaom-av1
				"-map 0:v:0 " +
				"-map 0:a? " +
				"-map 0:s? " +
				"-pix_fmt yuv420p " + // YUV 4:2:0
				"-profile:v high " + // force 8 bit
				fmt.Sprintf("-crf %d ", app.Setting.EncodingCrf) + // setting quality
				fmt.Sprintf("-filter:v scale=%d:-2 ", app.Setting.EncodingResolution) + // setting resolution
				fmt.Sprintf(`"%s" `, tmpOutput) +
				fmt.Sprintf("-progress unix://%s -y", tempSock(
					videoDuration,
					fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
					history,
				)) // progress tracking
	}
	startTime := time.Now()
	cmd := exec.Command(
		"bash",
		"-c",
		ffmpegCommand)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		log.Errorf("Error happend while encoding: %v\n", err.Error())
		log.Error("out", outb.String())
		log.Error("err", errb.String())
		log.Error(ffmpegCommand)

		if err := history.Failed(app.DB, fmt.Sprintf("%v | %v | %v", err.Error(), outb.String(), errb.String())); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}

	if err := history.Copy(app.DB, output); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

	if err := helper.Move(tmpOutput, output); err != nil {
		if err := history.Failed(app.DB, fmt.Sprintf("Failed to copy encoded file to output path: %v", err)); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}

	// delete original file
	if err := os.Remove(file); err != nil {
		log.Warn("Failed to delete old file\n", err)
	}
	// delete old nfo
	if err := os.Remove(fmt.Sprintf("%s.nfo", file)); err != nil {
		log.Warn("Failed to delete old nfo file\n", err)
	}

	endTime := time.Now()

	fi, err = os.Stat(output)
	if err != nil {
		log.Errorf("Failed to read filesize of new file %s\n", err)
		if err := history.Failed(app.DB, err.Error()); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}
	newSize := fi.Size()

	log.Infof("Old Size: %s / New Size: %s\n", humanize.Bytes(uint64(oldSize)), humanize.Bytes(uint64(newSize)))
	if err := history.Finished(app.DB, uint64(newSize), time.Duration(endTime.Unix()-startTime.Unix())*time.Second); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}
}

func tempSock(totalDuration float64, sockFileName string, encodingTask *m.History) string {
	sockFilePath := path.Join(os.TempDir(), sockFileName)
	l, err := net.Listen("unix", sockFilePath)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := ""
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = fmt.Sprintf("%.2f", float64(c)/totalDuration/1000000)
			}
			if strings.Contains(data, "progress=end") {
				cp = "1.0"
			}
			if cp == "" {
				cp = ".0"
			}
			if cp != progress {
				progress = cp
				// fmt.Println("progress: ", progress)
				floatProg, err := strconv.ParseFloat(progress, 64)
				if err != nil {
					fmt.Println("could not save progress in database")
				}
				if floatProg != 0 {
					encodingTask.SetProgress(app.DB, floatProg)
				}
			}
		}
	}()

	return sockFilePath
}
