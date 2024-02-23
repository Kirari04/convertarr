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
	Total   int64 // Total bytes to copy
	Copied  int64 // Bytes copied so far
	LastMsg time.Time
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.Writer.Write(p)
	pw.Copied += int64(n)
	progress := float64(pw.Copied) / float64(pw.Total) * 100
	if time.Since(pw.LastMsg).Seconds() > 2 {
		pw.LastMsg = time.Now()
		log.Infof("Copier in Progress... %.2f%%", progress)
	}
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
		Writer: dstFile,
		Total:  srcFileInfo.Size(),
		Copied: 0,
	}

	_, err = io.Copy(progressWriter, srcFile)
	if err != nil {
		return err
	}
	return nil
}
