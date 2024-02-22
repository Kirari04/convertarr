package helper

import (
	"bytes"
	"context"
	"encoder/app"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func Encode(inputFile, outputFile string) error {
	timeStart := time.Now()
	timeStartTotal := time.Now()
	// probe file so we can show encoding progress
	// ffprobe context
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	tmpDir := os.TempDir()
	tmpMetadataFile := fmt.Sprintf("%s/tmp.mkv", tmpDir)
	defer os.Remove(tmpMetadataFile)
	data, err := ffprobe.ProbeURL(ctx, inputFile)
	if err != nil {
		log.Errorf("Failed to update history %v", err)
		return err
	}

	log.Infof("Time Taken => FFProbe: %s", time.Since(timeStart))
	timeStart = time.Now()

	// export audio and subtitles
	ffmpegMetadata := fmt.Sprintf(`ffmpeg -i "%s" \
	-threads 0 \
	-map 0:a? \
	-map 0:s? \
	-c copy \
	"%s" -y`, inputFile, tmpMetadataFile)

	cmdM := exec.Command(
		"bash",
		"-c",
		ffmpegMetadata)

	var outbM, errbM bytes.Buffer
	cmdM.Stdout = &outbM
	cmdM.Stderr = &errbM

	if err := cmdM.Run(); err != nil {
		return fmt.Errorf(
			"error happend while exporting audios and subtitles: %v\nout: %v\nerr: %v\nCommand: %s",
			err.Error(), outbM.String(), errbM.String(), ffmpegMetadata,
		)
	}
	log.Infof("Time Taken => Export A/S: %s", time.Since(timeStart))
	timeStart = time.Now()

	videoDuration := data.Format.Duration().Seconds()
	m3u8Str := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-MEDIA-SEQUENCE:0"
	threadsPerEncode := 2
	allowThreads := app.Setting.EncodingThreads
	if allowThreads <= 0 {
		allowThreads = runtime.NumCPU()
	}
	threads := int(math.Floor(float64(allowThreads) / float64(threadsPerEncode)))
	threadChan := make(chan int, threads)
	size := app.Setting.EncodingResolution
	crf := app.Setting.EncodingCrf
	var chunckLen float64 = 15
	chuncks := int(math.Ceil(videoDuration / chunckLen))
	var wg sync.WaitGroup
	wg.Add(chuncks)
	var preSeekedChunck bool
	var encodingHasError error
	for i := 0; i < chuncks; i++ {
		if preSeekedChunck {
			log.Info("Preseeked chunck (skipping last)")
			wg.Done()
			continue
		}
		if encodingHasError != nil {
			wg.Done()
			continue
		}

		threadChan <- i
		tmpFile := fmt.Sprintf("%s/output%d.ts", tmpDir, i)
		from := chunckLen * (float64(i) + 0)
		to := chunckLen * (float64(i) + 1)
		nextFrom := chunckLen * (float64(i) + 1)
		nextTo := chunckLen * (float64(i) + 2)

		// change chunck size depeding on position in video
		if to > videoDuration {
			to = videoDuration
		}
		if nextTo > videoDuration {
			nextTo = videoDuration
		}
		if nextTo-nextFrom < chunckLen {
			to = videoDuration
			preSeekedChunck = true
		}
		tmpThreadsPerEncode := threadsPerEncode
		if preSeekedChunck {
			// allow last encodet to use more pools so encoding finishes faster
			tmpThreadsPerEncode = runtime.NumCPU()
		}

		m3u8Str += fmt.Sprintf("\n#EXTINF:%.6f,\noutput%d.ts", to-from, i)
		log.Infof("%s from: %.2f to: %.2f time: %s", tmpFile, from, to, time.Since(timeStart))
		defer os.Remove(tmpFile)
		var ffmpegCmd string
		if app.Setting.EnableHevcEncoding {
			if app.Setting.EnableAmdGpuEncoding {
				log.Info("Encoding Hevc using Vaapi interface")
				ffmpegCmd = fmt.Sprintf(`ffmpeg \
					-vaapi_device /dev/dri/renderD128 \
					-ss %.6f -t %.6f -i "%s" \
					-vf 'format=nv12,hwupload,scale_vaapi=-2:%d' \
					-threads %d \
					-c:v hevc_vaapi \
					-map 0:v:0 \
					-rc_mode CQP -pix_fmt vaapi -profile:v main \
					-global_quality %d \
					"%s" -y`,
					from,
					to-from,
					inputFile,
					size,
					tmpThreadsPerEncode,
					crf,
					tmpFile,
				)
			} else if app.Setting.EnableNvidiaGpuEncoding {
				log.Info("Encoding Hevc using Cuda interface")
				ffmpegCmd = fmt.Sprintf(`ffmpeg \
					-hwaccel_device 0 \
					-ss %.6f -t %.6f -i "%s" \
					-threads %d \
					-c:v hevc_nvenc \
					-map 0:v:0 \
					-rc:v vbr \
					-cq:v %d \
					-pix_fmt p010le -profile:v main \
					-vf "scale=%d:-2" \
					"%s" -y`,
					from,
					to-from,
					inputFile,
					tmpThreadsPerEncode,
					crf,
					size,
					tmpFile,
				)
			} else {
				log.Info("Encoding Hevc using Software interface")
				ffmpegCmd = fmt.Sprintf(`ffmpeg -ss %.6f -t %.6f -i "%s" \
					-threads 1 \
					-c:v libx265 \
					-map 0:v:0 \
					-pix_fmt yuv420p -profile:v main \
					-x265-params crf=%d:pools=%d -strict experimental \
					-filter:v scale=-2:%d \
					"%s" -y`,
					from,
					to-from,
					inputFile,
					crf, tmpThreadsPerEncode, size, tmpFile,
				)
			}
		} else {
			log.Info("Encoding H264 using Software interface")
			ffmpegCmd = fmt.Sprintf(`ffmpeg -ss %.6f -t %.6f -i "%s" \
					-threads %d \
					-c:v libx264 \
					-map 0:v:0 \
					-pix_fmt yuv420p -profile:v main \
					-crf %d \
					-filter:v scale=%d:-2 \
					"%s" -y`,
				from,
				to-from,
				inputFile,
				tmpThreadsPerEncode,
				crf, size, tmpFile,
			)
		}
		go func(i int) {
			defer wg.Done()
			defer func() {
				<-threadChan
			}()
			cmd := exec.Command(
				"bash",
				"-c",
				ffmpegCmd)

			var outb, errb bytes.Buffer
			cmd.Stdout = &outb
			cmd.Stderr = &errb

			if err := cmd.Run(); err != nil {
				encodingHasError = fmt.Errorf(
					"error happend while encoding chunck %d: %v\nout: %v\nerr: %v\nCommand: %s",
					i, err.Error(), outbM.String(), errbM.String(), ffmpegCmd,
				)
				return
			}
		}(i)
	}
	m3u8Str += "\n#EXT-X-ENDLIST"
	wg.Wait()

	// check if any chunck failed
	if encodingHasError != nil {
		return fmt.Errorf("encoding Failed with error: %v", encodingHasError)
	}

	log.Infof("Time Taken => Encode V: %s", time.Since(timeStart))
	timeStart = time.Now()

	// combine chuncks
	m3u8Path := fmt.Sprintf("%s/output.m3u8", tmpDir)
	defer os.Remove(m3u8Path)

	if err := os.WriteFile(m3u8Path, []byte(m3u8Str), 0644); err != nil {
		return fmt.Errorf("failed to create master.m3u8: %v", err)
	}

	ffmpegCmd := fmt.Sprintf(`ffmpeg -i "%s" \
	-i "%s" \
	-map 0:v:0 \
	-map 1:a? \
	-map 1:s? \
	-c copy \
	"%s" -y`, m3u8Path, tmpMetadataFile, outputFile)

	// create comparison image
	cmd := exec.Command(
		"bash",
		"-c",
		ffmpegCmd)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"error happend while combining video with audios and subtitles: %v\nout: %v\nerr: %v\nCommand: %s",
			err.Error(), outbM.String(), errbM.String(), ffmpegMetadata,
		)
	}
	log.Infof("Time Taken => Assemble A/S/V: %s", time.Since(timeStart))
	log.Infof("Time Taken => Total: %s", time.Since(timeStartTotal))

	return nil
}
