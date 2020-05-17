package input

import (
	"fmt"
	"github.com/bushaHQ/httputil/errors"
	"github.com/dchest/uniuri"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	MaxFileSize = 100 << 20
)

type link struct {
	dir string
}

func NewLink(dir string) Interface {
	return &link{dir: dir}
}

func (l link) Fetch(p string) (string, error) {

	s := time.Now()
	resp, err := http.Get(p)
	if err != nil {
		return "", errors.New("could not get video file", 401)
	}

	defer resp.Body.Close()

	cl := resp.Header.Get("Content-Length")
	log.Println("content-length", cl)
	clInt, err := strconv.Atoi(cl)

	//if err != nil {
	//	return "", errors.New("link does not specify file size")
	//}

	if err == nil && clInt > MaxFileSize {
		return "", errors.New(fmt.Sprintf("video is greater than %.4f MB", float64(MaxFileSize)/(1<<20)))
	}

	log.Println(resp.Header)

	outputPath := path.Join(l.dir, uniuri.NewLen(10))

	out, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.CopyN(out, resp.Body, MaxFileSize)
	if err != nil {
		if err != io.EOF {
			return "", err
		}
	}

	log.Println("GETTER TIME", time.Since(s))

	return outputPath, nil
}
