package helper

import (
	"io"
	"os"

	"github.com/labstack/gommon/log"
)

func Copy(srcpath, dstpath string) (err error) {
	r, err := os.Open(srcpath)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dstpath)
	if err != nil {
		return err
	}
	defer func() {
		if c := w.Close(); err == nil {
			err = c
		}
	}()

	// Create a multi writer to write to both the file and the logger
	multiWriter := io.MultiWriter(w, os.Stdout)

	// Create a proxy reader to intercept the data being read
	proxyReader := &ProxyReader{
		Reader:   r,
		Progress: logProgress,
	}

	_, err = io.Copy(multiWriter, proxyReader)
	return err
}

type ProxyReader struct {
	io.Reader
	Progress func(copied, total int64)
	Copied   int64
	Total    int64
}

func (pr *ProxyReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Copied += int64(n)
	pr.Progress(pr.Copied, pr.Total)
	return
}

func logProgress(copied, total int64) {
	log.Infof("Copying... %d/%d bytes", copied, total)
}
