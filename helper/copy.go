package helper

import (
	"io"
	"os"
	"time"

	"github.com/labstack/gommon/log"
)

// ProgressWriter wraps an io.Writer and prints progress information.
type ProgressWriter struct {
	io.Writer
	Total        int64 // Total bytes to copy
	Copied       int64 // Bytes copied so far
	SrcPath      string
	DstPath      string
	StopTracking bool
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.Writer.Write(p)
	pw.Copied += int64(n)
	return
}

func Copy(srcpath, dstpath string) (err error) {
	srcFile, err := os.Open(srcpath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dstpath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	progressWriter := &ProgressWriter{
		Writer:       dstFile,
		Total:        srcFileInfo.Size(),
		Copied:       0,
		SrcPath:      srcpath,
		DstPath:      dstpath,
		StopTracking: false,
	}

	go func() {
		for {
			time.Sleep(time.Second * 2)
			if progressWriter.StopTracking {
				break
			}
			progress := float64(progressWriter.Copied) / float64(progressWriter.Total) * 100
			log.Infof("Copier in Progress on file [%s] : %.2f%%", progressWriter.SrcPath, progress)
		}
	}()

	_, err = io.Copy(progressWriter, srcFile)
	if err != nil {
		return err
	}

	progressWriter.StopTracking = true

	return nil
}
