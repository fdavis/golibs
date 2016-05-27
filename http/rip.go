package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type RipSequence struct {
	DestinationDirectory string
	URLTemplate          string
	Min, Max             int
	Sleep                int
}

func GetSequence(r RipSequence) error {
	var url string
	if err := os.MkdirAll(r.DestinationDirectory, os.ModeDir|0755); err != nil {
		return fmt.Errorf("Fatal error: cannot create destination directory %s\n", r.DestinationDirectory)
	}
	for seqNum := r.Min; seqNum <= r.Max; seqNum++ {
		url = fmt.Sprintf(r.URLTemplate, seqNum)
		download(url, r.DestinationDirectory, seqNum)
		time.Sleep(time.Duration(r.Sleep) * time.Millisecond)
	}
	return nil
}

func download(url, dest string, num int) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Sprint(err)
		return
	}

	fileName := url[strings.LastIndex(url, "/")+1:]
	filePath := dest + "/" + strconv.Itoa(num) + "-" + fileName

	outf, err := os.Create(filePath)
	if err != nil {
		fmt.Sprint(err)
		return
	}

	_, err = io.Copy(outf, resp.Body)
	err = outf.Sync()

	resp.Body.Close()
	if err != nil {
		fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
}
