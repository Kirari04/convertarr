package setup

import (
	"bytes"
	"encoder/app"
	"encoder/m"
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
		// update old states of histories on startup
		var histories []m.History
		if err := app.DB.
			Not(&m.History{
				Status: "finished",
			}).
			Or(&m.History{
				Status: "failed",
			}).
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

	output := strings.TrimSuffix(file, ".mkv")
	output = fmt.Sprintf("%s[encoded]%s", output, ".mkv")

	if err := history.SetNewPath(app.DB, output); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

	if err := history.Encoding(app.DB); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

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

		if err := history.Failed(app.DB, fmt.Sprintf("%v | %v | %v", err.Error(), outb.String(), errb.String())); err != nil {
			log.Errorf("Failed to update history %v\n", err)
		}
		return
	}

	if err := history.Copy(app.DB, output); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}

	// delete original file
	if err := os.Remove(file); err != nil {
		log.Warn("Failed to delete old file\n", err)
	}
	// delete old nfo
	if err := os.Remove(fmt.Sprintf("%s.nfo", file)); err != nil {
		log.Warn("Failed to delete old nfo file\n", err)
	}

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
	if err := history.Finished(app.DB, uint64(oldSize), uint64(newSize)); err != nil {
		log.Errorf("Failed to update history %v\n", err)
	}
}
