package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type RipSequence struct {
	DestinationDirectory string
	URLTemplate          string
	Min, Max             int
}

func GetSequence(r RipSequence) error {
	ch := make(chan string)
	var url string
	if err := os.MkdirAll(r.DestinationDirectory, os.ModeDir|0755); err != nil {
		return fmt.Errorf("Fatal error: cannot create destination directory %s\n", r.DestinationDirectory)
	}
	for seqNum := r.Min; seqNum <= r.Max; seqNum++ {
		url = fmt.Sprintf(r.URLTemplate, seqNum)
		go download(url, r.DestinationDirectory, seqNum, ch)
	}
	for reqs := r.Min; reqs <= r.Max; reqs++ {
		_ = <-ch
	}
	return nil
}

func download(url, dest string, num int, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	fileName := url[strings.LastIndex(url, "/")+1:]
	filePath := dest + "/" + strconv.Itoa(num) + "-" + fileName

	outf, err := os.Create(filePath)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	_, err = io.Copy(outf, resp.Body)
	err = outf.Sync()

	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	ch <- ""
}
