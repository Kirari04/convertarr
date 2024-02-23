package setup

import (
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

func Scanner() {
	if app.Setting.AutomaticScannsAtStartup {
		go func() {
			ScannFolders()
		}()
	}
	go func() {
		for {
			time.Sleep(time.Second * 30)
			if app.Setting.LastFolderScann.Before(time.Now().Add(app.Setting.AutomaticScannsInterval * -1)) {
				ScannFolders()
			}
		}
	}()
}

func ScannFolders() {
	if app.IsFileScanning {
		log.Info("Already scanning folders")
		return
	}
	app.IsFileScanning = true
	defer func() {
		app.IsFileScanning = false
	}()
	log.Info("Starting scanning of folders")
	app.Setting.LastFolderScann = time.Now()
	app.Setting.Save(app.DB)
	var folders []m.Folder
	if err := app.DB.Find(&folders).Error; err != nil {
		log.Error("failed to list folders: ", err)
		return
	}
	var filesToEncode []string
	var nFiles uint64
	for _, rootFolder := range folders {
		files, addNFiles, err := helper.OSReadDir(rootFolder.Path, 0)
		if err != nil {
			log.Warnf("failed to walk path %s: %v", rootFolder.Path, err)
			continue
		}
		nFiles += addNFiles
		filesToEncode = append(filesToEncode, files...)
	}

	log.Infof("Found %d unencoded files", len(filesToEncode))
	var wg sync.WaitGroup
	ch := make(chan int, 1000)
	for _, fileToEncode := range filesToEncode {
		wg.Add(1)
		ch <- 0
		go func(fileToEncode string) {
			defer func() {
				<-ch
				wg.Done()
			}()
			var exists bool
			for _, existingFile := range app.FilesToEncode {
				if existingFile == fileToEncode {
					exists = true
				}
			}
			if exists {
				return
			}
			if app.Setting.EncodingMaxRetry > 0 {
				hash, err := helper.HashFile(fileToEncode)
				log.Debug("Failed to hash file to encode", err)
				if err != nil {
					return
				}
				var tries int64
				if err := app.DB.
					Model(&m.History{}).
					Where(&m.History{Hash: hash}).
					Count(&tries).Error; err != nil {
					log.Error("Failed to count encoding tries: ", err)
					return
				}
				if tries >= int64(app.Setting.EncodingMaxRetry) {
					log.Debug("Reached max retries of file ", fileToEncode)
					return
				}
			}
			app.FilesToEncode = append(app.FilesToEncode, fileToEncode)
		}(fileToEncode)
	}
	wg.Wait()
	now := time.Now()
	app.LastFileScan = &now
	app.LastScanNFiles = nFiles
	log.Infof("Finished Scanning %d files", nFiles)
}
