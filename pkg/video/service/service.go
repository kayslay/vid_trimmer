package service

import (
	"fmt"
	"github.com/bushaHQ/httputil/errors"
	"github.com/dchest/uniuri"
	"github.com/jtguibas/cinema"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kayslay/vid_trimmer/internal/input"
	"io"
	"net/http"
	"os"
	"path"
)

type Service interface {
	Download(w io.Writer, d DownloadStruct) (string, error)
}

type basicService struct {
	url input.Interface
	//youtube input.Interface
	//twitter input.Interface
}

func NewBasicService(url input.Interface) Service {
	return &basicService{url: url}
}

func (s basicService) Download(w io.Writer, d DownloadStruct) (string, error) {
	//	for now we only support url style downloads
	pathUrl, err := s.url.Fetch(d.URL)
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

	defer outputFile.Close()

	io.Copy(w, outputFile)

	return path.Base(outputPath), nil
}
