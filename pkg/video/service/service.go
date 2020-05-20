package service

import (
	"context"
	"fmt"
	"github.com/bushaHQ/httputil/errors"
	"github.com/dchest/uniuri"
	"github.com/jtguibas/cinema"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kayslay/vid_trimmer/internal/input"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Service interface {
	Download(ctx context.Context, w io.Writer, d DownloadStruct) (string, error)
}

type basicService struct {
	url     input.Interface
	youtube input.Interface
	//twitter input.Interface
}

func NewBasicService(url, youtube input.Interface) Service {
	return &basicService{url: url, youtube: youtube}
}

func (s basicService) Download(ctx context.Context, w io.Writer, d DownloadStruct) (string, error) {
	//	for now we only support url style downloads
	u, err := url.Parse(d.URL)
	if err != nil {
		return "", err
	}

	n := time.Now()
	pathUrl, err := s.getInput(u.Hostname()).Fetch(ctx, d.URL)
	defer input.Remove(pathUrl)

	if err != nil {
		return "", errors.CoverErr(
			err,
			errors.New("could not create file at the moment", http.StatusServiceUnavailable),
			log.WithFields(log.Fields{
				"context": "fetch",
				"method":  "video/download",
			}),
		)
	}
	log.Println("Time taken to process input", time.Since(n))

	n = time.Now()
	v, err := cinema.Load(pathUrl)
	if err != nil {
		return "", errors.CoverErr(
			err,
			errors.New("could not create file at the moment", http.StatusServiceUnavailable),
			log.WithFields(log.Fields{
				"context": "cinema/load",
				"method":  "video/download",
			}),
		)
	}

	outputPath := fmt.Sprintf("%s.%s", path.Join(os.TempDir(), uniuri.NewLen(10)), d.Type)
	defer input.Remove(outputPath)

	v.Trim(d.Start, d.End)
	v.Render(outputPath)

	outputFile, err := os.Open(outputPath)
	if err != nil {
		return "", errors.CoverErr(
			err,
			errors.New("could not create file at the moment", http.StatusServiceUnavailable),
			log.WithFields(log.Fields{
				"context": "cinema/load",
				"method":  "video/download",
			}),
		)
	}
	log.Println("Time cinema took to trim file", time.Since(n))

	n = time.Now()
	defer outputFile.Close()

	io.Copy(w, outputFile)
	log.Println("Time copy video", time.Since(n))

	return path.Base(outputPath), nil
}

func (s basicService) getInput(hostname string) input.Interface {
	switch {
	case strings.Contains(hostname, "youtu"):
		return s.youtube
	default:
		return s.url

	}
}
