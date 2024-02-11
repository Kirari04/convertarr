package helper

import (
	"encoder/app"
	"encoder/m"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func TempSock(totalDuration float64, sockFileName string, encodingTask *m.History) string {
	sockFilePath := path.Join(os.TempDir(), sockFileName)
	l, err := net.Listen("unix", sockFilePath)
	if err != nil {
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		progress := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := ""
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = fmt.Sprintf("%.2f", float64(c)/totalDuration/1000000)
			}
			if strings.Contains(data, "progress=end") {
				cp = "1.0"
			}
			if cp == "" {
				cp = ".0"
			}
			if cp != progress {
				progress = cp
				// fmt.Println("progress: ", progress)
				floatProg, err := strconv.ParseFloat(progress, 64)
				if err != nil {
					fmt.Println("could not save progress in database")
				}
				if floatProg != 0 {
					encodingTask.SetProgress(app.DB, floatProg)
				}
			}
		}
	}()

	return sockFilePath
}
