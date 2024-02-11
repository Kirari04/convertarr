package ffmpegcmd

import (
	"crypto/sha256"
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"fmt"

	"github.com/google/uuid"
)

func H264Cpu(file string, tmpOutput string, videoDuration float64, history *m.History) string {
	return "ffmpeg " +
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
		fmt.Sprintf("-progress unix://%s -y", helper.TempSock(
			videoDuration,
			fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
			history,
		)) // progress tracking
}
