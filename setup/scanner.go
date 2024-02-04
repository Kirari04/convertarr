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
			scannFolders()
		}()
	}
	go func() {
		for {
			time.Sleep(time.Second * 30)
			if app.Setting.LastFolderScann.Before(time.Now().Add(app.Setting.AutomaticScannsInterval * -1)) {
				scannFolders()
			}
		}
	}()
}

func scannFolders() {
	log.Info("Starting scanning of folders")
	app.Setting.LastFolderScann = time.Now()
	app.Setting.Save(app.DB)
	var folders []m.Folder
	if err := app.DB.Find(&folders).Error; err != nil {
		log.Error("failed to list folders: ", err)
		return
	}
	var filesToEncode []string
	for _, rootFolder := range folders {
		files, err := helper.OSReadDir(rootFolder.Path)
		if err != nil {
			log.Warnf("failed to walk path %s: %v", rootFolder.Path, err)
			continue
		}
		filesToEncode = append(filesToEncode, files...)
	}

	log.Infof("Found %d files for encoding", len(filesToEncode))
	for _, fileToEncode := range filesToEncode {
		var exists bool
		for _, existingFile := range app.FilesToEncode {
			if existingFile == fileToEncode {
				exists = true
			}
		}
		if !exists {
			app.FilesToEncode = append(app.FilesToEncode, fileToEncode)
		}
	}
}
