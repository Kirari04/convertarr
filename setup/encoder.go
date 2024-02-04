package setup

import (
	"bytes"
	"encoder/app"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/labstack/gommon/log"
)

func Encoder() {
	go func() {
		for {
			time.Sleep(time.Second * 5)
			if len(app.FilesToEncode) > 0 {
				fileToEncode := app.FilesToEncode[0]
				encodeFile(fileToEncode)
				app.FilesToEncode = app.FilesToEncode[1:]
			}
		}
	}()
}

func encodeFile(file string) {
	app.CurrentFileToEncode = file
	defer func() {
		app.CurrentFileToEncode = ""
	}()

	if strings.Contains(file, "[encoded]") {
		log.Infof("Skipping already encoded file %s\n", file)
		return
	}

	log.Infof("Encoding file %s\n", file)
	fi, err := os.Stat(file)
	if err != nil {
		log.Infof("Failed to read filesize %s\n", err)
		return
	}
	oldSize := fi.Size()

	output := strings.TrimSuffix(file, ".mkv")
	output = fmt.Sprintf("%s[encoded]%s", output, ".mkv")
	ffmpegCommand :=
		"ffmpeg " +
			fmt.Sprintf(`-i "%s" `, file) + // input file
			fmt.Sprintf("-threads %d ", app.Setting.EncodingThreads) +
			"-c:a copy " +
			"-c:s copy " +
			"-c:v libx264 " + // setting video codec libx264 | libaom-av1
			"-map 0 " +
			"-pix_fmt yuv420p " + // YUV 4:2:0
			"-profile:v high " + // force 8 bit
			fmt.Sprintf("-crf %d ", app.Setting.EncodingCrf) + // setting quality
			"-filter:v scale=1920:-1 " + // setting resolution
			"-y " +
			fmt.Sprintf(`"%s"`, output)

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
		return
	}
	// delete original file
	if err := os.Remove(file); err != nil {
		log.Warn("Failed to delete old file\n", err)
	}
	// delete nfo
	if err := os.Remove(fmt.Sprintf("%s.nfo", file)); err != nil {
		log.Warn("Failed to delete old nfo file\n", err)
	}

	fi, err = os.Stat(output)
	if err != nil {
		log.Errorf("Failed to read filesize of new file %s\n", err)
		return
	}
	newSize := fi.Size()

	log.Infof("Old Size: %s / New Size: %s\n", humanize.Bytes(uint64(oldSize)), humanize.Bytes(uint64(newSize)))
}
