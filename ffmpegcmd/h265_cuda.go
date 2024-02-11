package ffmpegcmd

import (
	"crypto/sha256"
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"fmt"

	"github.com/google/uuid"
)

func H265Cuda(file string, tmpOutput string, videoDuration float64, history *m.History) string {
	return "ffmpeg " +
		"-hwaccel_device 0 " +
		fmt.Sprintf(`-i "%s" `, file) + // input file
		fmt.Sprintf("-threads %d ", app.Setting.EncodingThreads) +
		"-c:v hevc_nvenc " + // setting video codec
		"-c:a copy " +
		"-c:s copy " +
		"-map 0:v:0 " +
		"-map 0:a? " +
		"-map 0:s? " +
		"-rc:v vbr " +
		fmt.Sprintf("-cq:v %d ", app.Setting.EncodingCrf) + // setting quality
		"-pix_fmt p010le " +
		"-profile:v main " +
		fmt.Sprintf(`-vf "scale=%d:-2" `, app.Setting.EncodingResolution) +
		fmt.Sprintf(`"%s" `, tmpOutput) +
		fmt.Sprintf("-progress unix://%s -y", helper.TempSock(
			videoDuration,
			fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
			history,
		)) // progress tracking
}
