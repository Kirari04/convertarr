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
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"gopkg.in/vansante/go-ffprobe.v2"
)

// IsEncoderAvailable checks if a specific ffmpeg encoder is available and if the
// corresponding hardware/drivers are detected on the system.
func IsEncoderAvailable(encoderName string) bool {
	// First, check if the encoder is listed in the ffmpeg build.
	cmd := exec.Command("ffmpeg", "-encoders")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Errorf("Failed to run `ffmpeg -encoders`: %v", err)
		return false // Assume not available if ffmpeg command fails
	}
	re, err := regexp.Compile(`\b` + encoderName + `\b`)
	if err != nil {
		log.Errorf("Failed to compile regex for encoder check: %v", err)
		return false
	}
	if !re.MatchString(out.String()) {
		// Encoder not found in the ffmpeg build, so it's definitely not available.
		log.Warnf("Encoder '%s' not found in `ffmpeg -encoders` output.", encoderName)
		return false
	}

	// Next, check for the presence of the required hardware/drivers.
	isNvidia := strings.Contains(encoderName, "nvenc")
	isAmd := strings.Contains(encoderName, "vaapi")

	if isNvidia {
		// Check for NVIDIA hardware by trying to run `nvidia-smi`.
		cmdNv := exec.Command("nvidia-smi")
		if errNv := cmdNv.Run(); errNv != nil {
			log.Warnf("NVIDIA encoder '%s' is in ffmpeg build, but 'nvidia-smi' command failed. Hardware/drivers may be missing or not installed correctly.", encoderName)
			return false
		}
	}

	if isAmd {
		// Check for AMD hardware by looking for a VAAPI render device node.
		files, errAmd := os.ReadDir("/dev/dri")
		if errAmd != nil {
			log.Warnf("Could not read /dev/dri to check for AMD GPU: %v", errAmd)
			return false
		}
		foundRenderNode := false
		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), "renderD") {
				foundRenderNode = true
				break
			}
		}
		if !foundRenderNode {
			log.Warnf("AMD encoder '%s' is in ffmpeg build, but no /dev/dri/renderD* node was found. Hardware/drivers may be missing or not installed correctly.", encoderName)
			return false
		}
	}

	// If both software and hardware checks pass, the encoder is available.
	return true
}

func Encode(inputFile, outputFile string, history *m.History) error {
	timeStart := time.Now()
	timeStartTotal := time.Now()
	// probe file so we can show encoding progress
	// ffprobe context
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, inputFile)
	if err != nil {
		log.Errorf("Failed to probe file %s: %v", inputFile, err)
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

	// generate ffmpeg command
	var ffmpegCmd string
	progressSocketArg := fmt.Sprintf("-progress unix://%s -y", TempSock(
		videoDuration,
		fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
		history,
	))

	if app.Setting.EnableHevcEncoding {
		// HEVC (H.265) Encoding Logic
		if app.Setting.EnableAmdGpuEncoding {
			log.Info("Encoding HEVC using AMD GPU (VAAPI)")
			ffmpegCmd = fmt.Sprintf(`ffmpeg \
                    -vaapi_device /dev/dri/renderD128 \
                    -i "%s" \
                    -vf 'scale=%d:-2,format=nv12,hwupload' \
                    -c:v hevc_vaapi \
                    -map 0:v:0 \
                    -c:a copy \
                    -c:s copy \
                    -map 0:a? \
                    -map 0:s? \
                    -rc_mode CQP -pix_fmt vaapi -profile:v main \
                    -global_quality %d \
                    "%s" %s`,
				inputFile,
				size,
				crf,
				outputFile,
				progressSocketArg,
			)
		} else if app.Setting.EnableNvidiaGpuEncoding {
			log.Info("Encoding HEVC using NVIDIA GPU (NVENC)")
			ffmpegCmd = fmt.Sprintf(`ffmpeg \
                    -hwaccel cuda -hwaccel_output_format cuda \
                    -i "%s" \
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
                    "%s" %s`,
				inputFile,
				crf,
				size,
				outputFile,
				progressSocketArg,
			)
		} else {
			log.Info("Encoding HEVC using CPU (libx265)")
			ffmpegCmd = fmt.Sprintf(`ffmpeg -i "%s" \
                    -c:v libx265 \
                    -map 0:v:0 \
                    -c:a copy \
                    -c:s copy \
                    -map 0:a? \
                    -map 0:s? \
                    -pix_fmt yuv420p -profile:v main \
                    -x265-params crf=%d:pools=%d -strict experimental \
                    -filter:v scale=%d:-2 \
                    "%s" %s`,
				inputFile,
				crf, allowThreads, size, outputFile,
				progressSocketArg,
			)
		}
	} else {
		// H.264 Encoding Logic
		if app.Setting.EnableAmdGpuEncoding {
			log.Info("Encoding H.264 using AMD GPU (VAAPI)")
			ffmpegCmd = fmt.Sprintf(`ffmpeg \
                    -vaapi_device /dev/dri/renderD128 \
                    -i "%s" \
                    -vf 'scale=%d:-2,format=nv12,hwupload' \
                    -c:v h264_vaapi \
                    -map 0:v:0 \
                    -c:a copy \
                    -c:s copy \
                    -map 0:a? \
                    -map 0:s? \
                    -rc_mode CQP -pix_fmt vaapi \
                    -global_quality %d \
                    "%s" %s`,
				inputFile,
				size,
				crf,
				outputFile,
				progressSocketArg,
			)
		} else if app.Setting.EnableNvidiaGpuEncoding {
			log.Info("Encoding H.264 using NVIDIA GPU (NVENC)")
			ffmpegCmd = fmt.Sprintf(`ffmpeg \
                    -hwaccel cuda -hwaccel_output_format cuda \
                    -i "%s" \
                    -c:v h264_nvenc \
                    -map 0:v:0 \
                    -c:a copy \
                    -c:s copy \
                    -map 0:a? \
                    -map 0:s? \
                    -rc:v vbr \
                    -cq:v %d \
                    -pix_fmt yuv420p -profile:v main \
                    -vf "scale=%d:-2" \
                    "%s" %s`,
				inputFile,
				crf,
				size,
				outputFile,
				progressSocketArg,
			)
		} else {
			log.Info("Encoding H.264 using CPU (libx264)")
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
                    "%s" %s`,
				inputFile,
				allowThreads,
				crf, size, outputFile,
				progressSocketArg,
			)
		}
	}

	cmd := exec.Command(
		"bash",
		"-c",
		ffmpegCmd)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	log.Infof("Executing FFMpeg command: %s", cmd.String())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"error happened while encoding: %v\nout: %v\nerr: %v\nCommand: %s",
			err.Error(), outb.String(), errb.String(), ffmpegCmd,
		)
	}

	log.Infof("Time Taken => FFMpeg Encode: %s", time.Since(timeStart))
	log.Infof("Time Taken => Total: %s", time.Since(timeStartTotal))

	return nil
}
