package helper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoder/app"
	"encoder/m"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func Encode(inputFile, outputFile string, history *m.History) error {
	timeStart := time.Now()
	timeStartTotal := time.Now()
	// probe file so we can show encoding progress
	// ffprobe context
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	tmpDir := os.TempDir()
	tmpMetadataFile := fmt.Sprintf("%s/tmpmetadata.mkv", tmpDir)
	defer os.Remove(tmpMetadataFile)
	data, err := ffprobe.ProbeURL(ctx, inputFile)
	if err != nil {
		log.Errorf("Failed to update history %v", err)
		return err
	}

	log.Infof("Time Taken => FFProbe: %s", time.Since(timeStart))
	timeStart = time.Now()

	videoDuration := data.Format.Duration().Seconds()

	allowThreads := app.Setting.EncodingThreads
	if allowThreads <= 0 {
		allowThreads = runtime.NumCPU()
	}
	size := app.Setting.EncodingResolution
	crf := app.Setting.EncodingCrf

	// generate ffmpeg code
	var ffmpegCmd string
	if app.Setting.EnableHevcEncoding {
		if app.Setting.EnableAmdGpuEncoding {
			log.Info("Encoding Hevc using Vaapi interface")
			ffmpegCmd = fmt.Sprintf(`ffmpeg \
					-vaapi_device /dev/dri/renderD128 \
					-i "%s" \
					-vf 'format=nv12,hwupload,scale_vaapi=%d:-2' \
					-threads %d \
					-c:v hevc_vaapi \
					-map 0:v:0 \
					-c:a copy \
					-c:s copy \
					-map 0:a? \
					-map 0:s? \
					-rc_mode CQP -pix_fmt vaapi -profile:v main \
					-global_quality %d \
					"%s" %s -y`,
				inputFile,
				size,
				allowThreads,
				crf,
				outputFile,
				fmt.Sprintf("-progress unix://%s -y", TempSock(
					videoDuration,
					fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
					history,
				)), // progress tracking
			)
		} else if app.Setting.EnableNvidiaGpuEncoding {
			log.Info("Encoding Hevc using Cuda interface")
			ffmpegCmd = fmt.Sprintf(`ffmpeg \
					-hwaccel_device 0 \
					-i "%s" \
					-threads %d \
					-c:v hevc_nvenc \
					-map 0:v:0 \
					-c:a copy \
					-c:s copy \
					-map 0:a? \
					-map 0:s? \
					-rc:v vbr \
					-cq:v %d \
					-pix_fmt p010le -profile:v main \
					-vf "scale=%d:-2" \
					"%s" %s -y`,
				inputFile,
				allowThreads,
				crf,
				size,
				outputFile,
				fmt.Sprintf("-progress unix://%s -y", TempSock(
					videoDuration,
					fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
					history,
				)), // progress tracking
			)
		} else {
			log.Info("Encoding Hevc using Software interface")
			ffmpegCmd = fmt.Sprintf(`ffmpeg -i "%s" \
					-threads 1 \
					-c:v libx265 \
					-map 0:v:0 \
					-c:a copy \
					-c:s copy \
					-map 0:a? \
					-map 0:s? \
					-pix_fmt yuv420p -profile:v main \
					-x265-params crf=%d:pools=%d -strict experimental \
					-filter:v scale=%d:-2 \
					"%s" %s -y`,
				inputFile,
				crf, allowThreads, size, outputFile,
				fmt.Sprintf("-progress unix://%s -y", TempSock(
					videoDuration,
					fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
					history,
				)), // progress tracking
			)
		}
	} else {
		log.Info("Encoding H264 using Software interface")
		ffmpegCmd = fmt.Sprintf(`ffmpeg -i "%s" \
					-threads %d \
					-c:v libx264 \
					-map 0:v:0 \
					-c:a copy \
					-c:s copy \
					-map 0:a? \
					-map 0:s? \
					-pix_fmt yuv420p -profile:v main \
					-crf %d \
					-filter:v scale=%d:-2 \
					"%s" %s -y`,
			inputFile,
			allowThreads,
			crf, size, outputFile,
			fmt.Sprintf("-progress unix://%s -y", TempSock(
				videoDuration,
				fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
				history,
			)), // progress tracking
		)
	}
	cmd := exec.Command(
		"bash",
		"-c",
		ffmpegCmd)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"error happend while encoding: %v\nout: %v\nerr: %v\nCommand: %s",
			err.Error(), outb.String(), errb.String(), ffmpegCmd,
		)
	}

	log.Infof("Time Taken => Assemble A/S/V: %s", time.Since(timeStart))
	log.Infof("Time Taken => Total: %s", time.Since(timeStartTotal))

	return nil
}
