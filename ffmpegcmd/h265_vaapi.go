package ffmpegcmd

import (
	"crypto/sha256"
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"fmt"

	"github.com/google/uuid"
)

func H265Vaapi(file string, tmpOutput string, videoDuration float64, history *m.History) string {
	return "ffmpeg " +
		"-vaapi_device /dev/dri/renderD128 " +
		fmt.Sprintf(`-i "%s" `, file) + // input file
		fmt.Sprintf("-vf 'format=nv12,hwupload,scale_vaapi=%d:-2' ", app.Setting.EncodingResolution) +
		fmt.Sprintf("-threads %d ", app.Setting.EncodingThreads) +
		"-c:v hevc_vaapi " + // setting video codec
		"-c:a copy " +
		"-c:s copy " +
		"-map 0:v:0 " +
		"-map 0:a? " +
		"-map 0:s? " +
		"-rc_mode CQP " +
		"-pix_fmt vaapi_vld " +
		"-profile:v main " +
		fmt.Sprintf("-global_quality %d ", app.Setting.EncodingCrf) + // setting quality
		fmt.Sprintf(`"%s" `, tmpOutput) +
		fmt.Sprintf("-progress unix://%s -y", helper.TempSock(
			videoDuration,
			fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
			history,
		)) // progress tracking
}
