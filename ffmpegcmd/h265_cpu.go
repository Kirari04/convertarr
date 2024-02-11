package ffmpegcmd

import (
	"crypto/sha256"
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"fmt"

	"github.com/google/uuid"
)

func H265Cpu(file string, tmpOutput string, videoDuration float64, history *m.History) string {
	h265Pools := "*"
	if app.Setting.EncodingThreads > 0 {
		h265Pools = fmt.Sprint(app.Setting.EncodingThreads)
	}
	return "ffmpeg " +
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
		fmt.Sprintf("-progress unix://%s -y", helper.TempSock(
			videoDuration,
			fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
			history,
		)) // progress tracking
}
