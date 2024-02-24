package setup

import (
	"bytes"
	"context"
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"fmt"
	"os"
	"os/exec"
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
				log.Info("Encoding is disabled")
				continue
			}
			if len(app.FilesToEncode) > 0 {
				fileToEncode := app.FilesToEncode[0]
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

	output := strings.TrimSuffix(file, ".mkv")
	output = fmt.Sprintf("%s[encoded]%s", output, ".mkv")
	tmpDir := os.TempDir()
	tmpOutput := fmt.Sprintf("%s/tmp.mkv", tmpDir)
	tmpInput := fmt.Sprintf("%s/tmpin.mkv", tmpDir)

	defer os.Remove(tmpOutput)
	defer os.Remove(tmpInput)

	if app.Setting.PreCopyFileCount > 0 {
		log.Infof("Waiting on copier for file %s", file)
		// tell copier what file is asked
		app.AwaitForFileCopy = file
		// wait for copier to respons with the tmp file
		setTmpInput, ok := <-app.AwaitForFileCopyChan
		if !ok {
			if err := history.Failed(app.DB, "Failed to AwaitForFileCopyChan"); err != nil {
				log.Errorf("Failed to update history %v\n", err)
			}
			return
		}
		if err := helper.Move(setTmpInput, tmpInput); err != nil {
			os.Remove(setTmpInput)
			if err := history.Failed(app.DB, "Failed to move tmp file from copier"); err != nil {
				log.Errorf("Failed to update history %v\n", err)
			}
			return
		}
		log.Infof("Received file from copier %s", file)
		app.PreloadedFiles.Remove(file)
	} else {
		log.Infof("Copy file for Encoding to local folder %s\n", file)
		// copy file to locale path
		if err := helper.Copy(file, tmpInput); err != nil {
			if err := history.Failed(app.DB, fmt.Sprintf("Failed to copy encoded file to input path: %v", err)); err != nil {
				log.Errorf("Failed to update history %v\n", err)
			}
			return
		}
	}

	log.Infof("Encoding file %s\n", file)
	fi, err := os.Stat(tmpInput)
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

	hash, err := helper.HashFile(tmpInput)
	if err != nil {
		if err := history.Failed(app.DB, err.Error()); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}
	if err := history.SetHash(app.DB, hash); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

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
	data, err := ffprobe.ProbeURL(ctx, tmpInput)
	if err != nil {
		if err := history.Failed(app.DB, err.Error()); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}
	dataStreams := data.StreamType(ffprobe.StreamAny)
	// videoDuration := data.Format.Duration().Seconds()
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
	startTime := time.Now()
	if err := helper.Encode(tmpInput, tmpOutput, history); err != nil {
		if err := history.Failed(app.DB, err.Error()); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}

	// generate image comparison
	if app.Setting.EnableImageComparison {
		imgOutputPath := fmt.Sprintf("./imgs/%s.jpg", uuid.NewString())
		ffmpegImgCommand := `ffmpeg ` +
			`-t 1 -s 1920x2160 -f rawvideo -pix_fmt rgb24 -r 25 ` +
			`-i /dev/zero ` +
			fmt.Sprintf(`-i "%s" `, tmpInput) +
			fmt.Sprintf(`-i "%s" `, tmpOutput) +
			`-filter_complex ` +
			`"[0:v]scale=-2:2160[bg]; ` +
			`[1:v:0]scale=-2:1080[img]; ` +
			`[2:v:0]scale=-2:1080[img2]; ` +
			`[img]crop=iw/8:ih/8,scale=8*iw:-2[imgz]; ` +
			`[img2]crop=iw/8:ih/8,scale=8*iw:-2[img2z]; ` +
			`[imgz]split=1[v1]; ` +
			`[img2z]split=1[v2]; ` +
			`[bg][v1]overlay=w*0:h*0,trim=start=5:end=6[f2]; ` +
			`[f2][v2]overlay=w*0:h*1,trim=start=5:end=6[f3]; ` +
			`[f3]setpts=PTS-STARTPTS,scale=-2:2160[fin]" ` +
			`-map [fin] ` +
			`-qscale:v 2 ` +
			fmt.Sprintf(`-frames:v 1 "%s" -y`, imgOutputPath)

		// create comparison image
		cmdImg := exec.Command(
			"bash",
			"-c",
			ffmpegImgCommand)

		var cmdImgOutb, cmdImgErrb bytes.Buffer
		cmdImg.Stdout = &cmdImgOutb
		cmdImg.Stderr = &cmdImgErrb

		err := cmdImg.Run()
		if err != nil {
			log.Errorf("Error happend while creating comparison img: %v\n", err.Error())
			log.Error("out: ", cmdImgOutb.String())
			log.Error("err: ", cmdImgErrb.String())
			log.Error(ffmpegImgCommand)
		}
		if err == nil {
			if err := history.SetComparisonImg(app.DB, imgOutputPath); err != nil {
				log.Errorf("Failed to update history %v\n", err)
			}
		}
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
