package setup

import (
	"encoder/app"
	"encoder/helper"
	"encoder/m"
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
	app.IsFileScanning = true
	defer func() {
		app.IsFileScanning = false
	}()
	if app.IsFileScanning {
		log.Info("Already scanning folders")
		return
	}
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
	for _, fileToEncode := range filesToEncode {
		var exists bool
		for _, existingFile := range app.FilesToEncode {
			if existingFile == fileToEncode {
				exists = true
			}
		}
		if exists {
			continue
		}
		if app.Setting.EncodingMaxRetry > 0 {
			hash, err := helper.HashFile(fileToEncode)
			log.Debug("Failed to hash file to encode", err)
			if err != nil {
				continue
			}
			var tries int64
			if err := app.DB.
				Model(&m.History{}).
				Where(&m.History{Hash: hash}).
				Count(&tries).Error; err != nil {
				log.Error("Failed to count encoding tries: ", err)
				continue
			}
			if tries >= int64(app.Setting.EncodingMaxRetry) {
				log.Debug("Reached max retries of file ", fileToEncode)
				continue
			}
		}
		app.FilesToEncode = append(app.FilesToEncode, fileToEncode)
	}
	now := time.Now()
	app.LastFileScan = &now
	app.LastScanNFiles = nFiles
}
