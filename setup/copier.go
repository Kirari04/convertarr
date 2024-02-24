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
		if err := os.MkdirAll(tmpPath, 0777); err != nil {
			log.Errorf("Failed to create tmp folder: %v", err)
		}

		// listening for encoders requirements
		go func() {
			var encoderWaitingOnEmptyPreload int
			for {
				time.Sleep(time.Second)
				if app.AwaitForFileCopy != "" {
					log.Infof("Requiested copier for file: %s", app.AwaitForFileCopy)
					var foundRequiredFile bool
					for _, preloadedFile := range app.PreloadedFiles {
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
							app.PreloadedFiles = removePreloadedFile(preloadedFile.File, app.PreloadedFiles)
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
				(len(app.PreloadedFiles) < app.Setting.PreCopyFileCount || overwriteNextPreloadFile != "") && // max preloaded files check, skip if overwriteNextPreloadFile is set
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
						if hasPreloadedFile(fileToEncode, app.PreloadedFiles) {
							nthFileToEncode++
							continue
						}
					} else {
						// if for some reason the preloaders preloading predictions where wrong
						// we have to choose the next file manually because the encoder is waiting for it
						// else the encoder would get stuck for ever
						fileToEncode = overwriteNextPreloadFile
						if hasPreloadedFile(fileToEncode, app.PreloadedFiles) {
							nthFileToEncode++
							continue
						}
					}

					log.Infof("Starting copier on file: %s", fileToEncode)
					tmpFilePath := fmt.Sprintf("%s/%s.mkv", tmpPath, uuid.NewString())

					preloadedFile := t.PreloadedFile{
						File:    fileToEncode,
						TmpPath: tmpFilePath,
						IsReady: false,
					}
					app.PreloadedFiles = append(app.PreloadedFiles, &preloadedFile)
					overwriteNextPreloadFile = ""

					// copy file to tmp path
					if err := helper.Copy(fileToEncode, tmpFilePath); err != nil {
						os.Remove(tmpFilePath)
						log.Errorf("Copier Failed to copy file to tmp folder: %v", err)
						app.PreloadedFiles = removePreloadedFile(fileToEncode, app.PreloadedFiles)
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

func removePreloadedFile(a string, list []*t.PreloadedFile) []*t.PreloadedFile {
	var i = -1
	for ii, b := range list {
		if b.File == a {
			i = ii
			break
		}
	}
	if i == -1 {
		return list
	}
	// replace "to be deleted" with last element
	list[i] = list[len(list)-1]
	// return while array excluding the last element (that now sits on the to be replaced index)
	return list[:len(list)-1]
}

func hasPreloadedFile(a string, list []*t.PreloadedFile) bool {
	for _, b := range list {
		if b.File == a {
			return true
		}
	}
	return false
}
