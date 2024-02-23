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
			for {
				time.Sleep(time.Second)
				if app.AwaitForFileCopy != "" {
					log.Infof("Requiested copier for file: %s", app.AwaitForFileCopy)
					for _, preloadedFile := range app.PreloadedFiles {
						if preloadedFile.File == app.AwaitForFileCopy {
							app.AwaitForFileCopyChan <- preloadedFile.TmpPath
							app.AwaitForFileCopy = ""
							log.Infof("Copier responded with file: %s", preloadedFile.TmpPath)
							app.PreloadedFiles = removePreloadedFile(preloadedFile.File, app.PreloadedFiles)
						}
					}
				}
			}
		}()

		// predict needed files
		for {
			time.Sleep(2 * time.Second)
			nthFileToEncode := 0
			if app.Setting.EnableEncoding &&
				app.Setting.PreCopyFileCount > 0 &&
				len(app.PreloadedFiles) < app.Setting.PreCopyFileCount {
				if len(app.FilesToEncode) > 0 {
					for {
						if nthFileToEncode >= len(app.FilesToEncode) {
							// all possible files had been already preloaded
							break
						}
						// get file that hasn't been preloaded yet
						fileToEncode := app.FilesToEncode[nthFileToEncode]
						if hasPreloadedFile(fileToEncode, app.PreloadedFiles) {
							nthFileToEncode++
							continue
						}

						log.Infof("Starting copier on file: %s", fileToEncode)
						tmpFilePath := fmt.Sprintf("%s/%s.mkv", tmpPath, uuid.NewString())
						// copy file to tmp path first
						if err := helper.Copy(fileToEncode, tmpFilePath); err != nil {
							os.Remove(tmpFilePath)
							log.Errorf("Copier Failed to copy file to tmp folder: %v", err)
							app.PreloadedFiles = removePreloadedFile(fileToEncode, app.PreloadedFiles)
							continue
						}
						// add file to array (when its ready / finished copy) so we can find it when the encoder requires it
						preloadedFile := t.PreloadedFile{
							File:    fileToEncode,
							TmpPath: tmpFilePath,
						}
						app.PreloadedFiles = append(app.PreloadedFiles, preloadedFile)
						log.Infof("Finished copier on file: %s", fileToEncode)
					}
				}
			}
		}
	}()
}

func removePreloadedFile(a string, list []t.PreloadedFile) []t.PreloadedFile {
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

func hasPreloadedFile(a string, list []t.PreloadedFile) bool {
	for _, b := range list {
		if b.File == a {
			return true
		}
	}
	return false
}
