package helper

import (
	"fmt"
	"os"
	"strings"
)

func OSReadDir(root string, initNFiles uint64) (files []string, nFiles uint64, err error) {
	nFiles = initNFiles
	f, err := os.Open(root)
	if err != nil {
		return files, nFiles, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, nFiles, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			newRootDir := fmt.Sprintf("%s%s/", root, file.Name())
			if !strings.HasPrefix(file.Name(), "/") {
				newRootDir = fmt.Sprintf("%s/%s/", root, file.Name())
			}
			newFiles, addNFiles, err := OSReadDir(newRootDir, nFiles)
			nFiles += addNFiles
			if err != nil {
				return files, nFiles, err
			}
			files = append(files, newFiles...)
		} else {
			nFiles++
			// TODO: legacy cheeck for encoded in filename
			if !strings.Contains(file.Name(), "[encoded]") {
				if strings.HasSuffix(file.Name(), ".mkv") {

					newFilePath := fmt.Sprintf("%s%s", root, file.Name())
					if !strings.HasPrefix(file.Name(), "/") {
						newFilePath = fmt.Sprintf("%s/%s", root, file.Name())
					}
					files = append(files, newFilePath)
				}
			}

		}
	}
	return files, nFiles, nil
}
