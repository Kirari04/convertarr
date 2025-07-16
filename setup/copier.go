package setup

import (
	"encoder/app"
	"encoder/helper"
	"encoder/t"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

var overwriteNextPreloadFile string

func Copier() {
	go func() {
		tmpRootPath := fmt.Sprintf("%s/copier", os.TempDir())
		os.RemoveAll(tmpRootPath) // remove old tmp folders
		// create new tmp flder
		tmpPath := fmt.Sprintf("%s/%s", tmpRootPath, uuid.NewString())
		if err := os.MkdirAll(tmpPath, os.ModePerm); err != nil {
			log.Errorf("Failed to create tmp folder: %v", err)
		}

		// listening for encoders requirements
		go func() {
			var encoderWaitingOnEmptyPreload int
			for {
				time.Sleep(time.Second)
				if app.AwaitForFileCopy != "" {
					log.Infof("Copier searching for file: %s", app.AwaitForFileCopy)
					var foundRequiredFile bool
					for _, preloadedFile := range app.PreloadedFiles.Get() {
						if preloadedFile == nil {
							continue
						}
						if preloadedFile.File == app.AwaitForFileCopy {
							foundRequiredFile = true
							if !preloadedFile.IsReady {
								continue
							}
							app.AwaitForFileCopyChan <- preloadedFile.TmpPath
							app.AwaitForFileCopy = ""
							log.Infof("Copier responded with file: %s", preloadedFile.TmpPath)
						}
					}
					if !foundRequiredFile {
						encoderWaitingOnEmptyPreload++
						if encoderWaitingOnEmptyPreload > 30 {
							overwriteNextPreloadFile = app.AwaitForFileCopy
						}
					}
				}
			}
		}()

		// predict needed files
		for {
			time.Sleep(2 * time.Second)
			if app.Setting.EnableEncoding && // encoding is enabled
				app.Setting.PreCopyFileCount > 0 && // preload is enabled
				(len(app.PreloadedFiles.Get()) < app.Setting.PreCopyFileCount || overwriteNextPreloadFile != "") && // max preloaded files check, skip if overwriteNextPreloadFile is set
				!app.IsFileScanning && // not scanning => meaning app.FilesToEncode won't be altered
				len(app.FilesToEncode) > 0 { // has any files to encode
				// we have a for loop here so we can instantly choose another file if it already is preloaded
				nthFileToEncode := 0
				for {
					if nthFileToEncode >= len(app.FilesToEncode) {
						// all possible files had been already preloaded
						break
					}
					// get file that hasn't been preloaded yet
					var fileToEncode string
					if overwriteNextPreloadFile == "" {
						// default behavior
						fileToEncode = app.FilesToEncode[nthFileToEncode]
					} else {
						// if for some reason the preloaders preloading predictions where wrong
						// we have to choose the next file manually because the encoder is waiting for it
						// else the encoder would get stuck for ever
						fileToEncode = overwriteNextPreloadFile
						overwriteNextPreloadFile = ""
					}

					if app.CurrentFileToEncode == fileToEncode {
						nthFileToEncode++
						continue
					}

					if app.PreloadedFiles.Exists(fileToEncode) {
						nthFileToEncode++
						continue
					}

					log.Infof("Starting copier on file: %s", fileToEncode)
					tmpFilePath := fmt.Sprintf("%s/%s.mkv", tmpPath, uuid.NewString())
					preloadedFile := t.PreloadedFile{
						File:    fileToEncode,
						TmpPath: tmpFilePath,
						IsReady: false,
					}
					app.PreloadedFiles.Append(&preloadedFile)

					// copy file to tmp path
					if err := helper.Copy(fileToEncode, tmpFilePath); err != nil {
						os.Remove(tmpFilePath)
						log.Errorf("Copier Failed to copy file to tmp folder: %v", err)
						app.PreloadedFiles.Remove(fileToEncode)
						break
					}

					preloadedFile.IsReady = true

					log.Infof("Finished copier on file: %s", fileToEncode)
					break
				}
			}
		}
	}()
}
