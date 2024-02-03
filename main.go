package main

import (
	"bytes"
	"encoder/server"
	"encoder/setup"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/thatisuday/commando"
)

var (
	ROOT_DIR     = "./"
	VIDEO_SUFFIX = ".mkv"
)

func main() {
	commando.
		SetExecutableName("encoder").
		SetVersion("v1.0.0").
		SetDescription("This CLI tool encodes all mkv files.")

	commando.
		Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			fmt.Println("use help command")
		})

	commando.
		Register("serve").
		SetShortDescription("start webserver on :8080").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			setup.Setup()
			server.Serve()
		})

	commando.
		Register("list").
		SetShortDescription("list files").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

			files, err := OSReadDir(ROOT_DIR)
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				fmt.Println(file)
			}

			fmt.Printf("Found %d files\n", len(files))
		})

	commando.
		Register("encode").
		SetShortDescription("starts encoding").
		AddArgument("threads", "ffmpeg threads ussage", "0").
		AddArgument("crf", "ffmpeg crf ussage", "25").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			// validation
			threadsRaw := args["threads"].Value
			threads, err := strconv.Atoi(threadsRaw)
			if err != nil {
				log.Fatal(err)
			}
			if threads > runtime.NumCPU() {
				log.Fatalf("You can only use %d threads", runtime.NumCPU())
			}
			if threads < 0 {
				log.Fatal("Threads minimum is 0")
			}
			crfRaw := args["crf"].Value
			crf, err := strconv.Atoi(crfRaw)
			if err != nil {
				log.Fatal(err)
			}
			if crf < 0 || crf > 50 {
				log.Fatal("Crf out of range (0-50)")
			}

			// reead root directory
			files, err := OSReadDir(ROOT_DIR)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Found %d files\n", len(files))
			for _, file := range files {
				if strings.Contains(file, "[encoded]") {
					fmt.Printf("Skipping already encoded file %s\n", file)
					continue
				}
				fmt.Printf("Encoding file %s\n", file)
				fi, err := os.Stat(file)
				if err != nil {
					fmt.Printf("Failed to read filesize %s\n", err)
					continue
				}
				oldSize := fi.Size()

				output := strings.TrimSuffix(file, VIDEO_SUFFIX)
				output = fmt.Sprintf("%s[encoded]%s", output, VIDEO_SUFFIX)
				ffmpegCommand :=
					"ffmpeg " +
						fmt.Sprintf(`-i "%s" `, file) + // input file
						fmt.Sprintf("-threads %d ", threads) +
						"-c:a copy " +
						"-c:s copy " +
						"-c:v libx264 " + // setting video codec libx264 | libaom-av1
						"-map 0 " +
						"-pix_fmt yuv420p " + // YUV 4:2:0
						"-profile:v high " + // force 8 bit
						fmt.Sprintf("-crf %d ", crf) + // setting quality
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
					log.Printf("Error happend while encoding: %v\n", err.Error())
					log.Println("out", outb.String())
					log.Println("err", errb.String())
					log.Println(ffmpegCommand)
					time.Sleep(time.Second * 2)
					continue
				}
				// delete original file
				if err := os.Remove(file); err != nil {
					log.Println("Failed to delete old file\n", err)
				}
				// delete nfo
				if err := os.Remove(fmt.Sprintf("%s.nfo", file)); err != nil {
					log.Println("Failed to delete nfo file\n", err)
				}

				fi, err = os.Stat(output)
				if err != nil {
					fmt.Printf("Failed to read filesize %s\n", err)
					continue
				}
				newSize := fi.Size()

				fmt.Printf("Old Size: %s / New Size: %s\n", humanize.Bytes(uint64(oldSize)), humanize.Bytes(uint64(newSize)))

				// if err := os.Rename(output, file); err != nil {
				// 	log.Println("Failed to rename encodet file ", err)
				// }
			}
		})

	commando.Parse(nil)
}

func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {

		if file.IsDir() {
			newFiles, err := OSReadDir(fmt.Sprintf("%s%s/", root, file.Name()))
			if err != nil {
				return files, err
			}
			files = append(files, newFiles...)
		} else {
			if strings.HasSuffix(file.Name(), VIDEO_SUFFIX) {
				files = append(files, fmt.Sprintf("%s%s", root, file.Name()))
			}
		}
	}
	return files, nil
}
