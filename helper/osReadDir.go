package helper

import (
	"fmt"
	"os"
	"strings"
)

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
			newRootDir := fmt.Sprintf("%s%s/", root, file.Name())
			if !strings.HasPrefix(file.Name(), "/") {
				newRootDir = fmt.Sprintf("%s/%s/", root, file.Name())
			}
			newFiles, err := OSReadDir(newRootDir)
			if err != nil {
				return files, err
			}
			files = append(files, newFiles...)
		} else {
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
	return files, nil
}
