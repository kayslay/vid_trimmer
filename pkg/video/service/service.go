package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/bushaHQ/httputil/errors"
	"github.com/jtguibas/cinema"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kayslay/vid_trimmer/internal/filestore"
	"gitlab.com/kayslay/vid_trimmer/internal/input"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service interface {
	Download(ctx context.Context, d DownloadStruct) (string, bool, error)
	GeneratePath(key string) string
	Exists(key string) bool
}

type basicService struct {
	url       input.Interface
	youtube   input.Interface
	fileStore filestore.FileStorer
	//twitter input.Interface
}

func NewBasicService(url, youtube input.Interface, storer filestore.FileStorer) Service {
	return &basicService{url: url, youtube: youtube, fileStore: storer}
}

func (s basicService) Download(ctx context.Context, d DownloadStruct) (string, bool, error) {
	key, exist := s.generateFileKey(d)
	if exist {
		return key, exist, nil
	}

	outputPath := s.fileStore.GeneratePath(key)
	//	for now we only support url style downloads
	u, err := url.Parse(d.URL)
	if err != nil {
		return "", false, err
	}

	n := time.Now()
	pathUrl, err := s.getInput(u.Hostname()).Fetch(ctx, d.URL)

	if err != nil {
		return "", false, errors.CoverErr(
			err,
			errors.New("could not create file at the moment", http.StatusServiceUnavailable),
			log.WithFields(log.Fields{
				"context": "fetch",
				"method":  "video/download",
			}),
		)
	}
	log.Println("Time taken to process input", time.Since(n))

	//trim in the background
	go func() {
		defer input.Remove(pathUrl)
		defer func() {
			//handle panic
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"context": "cinema/load",
					"method":  "video/download",
				})
			}
		}()

		n := time.Now()
		v, err := cinema.Load(pathUrl)
		if err != nil {
			errors.CoverErr(
				err,
				errors.New("could not create file at the moment", http.StatusServiceUnavailable),
				log.WithFields(log.Fields{
					"context": "cinema/load",
					"method":  "video/download",
				}),
			)
			return
		}

		v.Trim(d.Start, d.End)
		err = v.Render(outputPath)

		if err != nil {
			errors.CoverErr(
				err,
				errors.New("could not create file at the moment", http.StatusServiceUnavailable),
				log.WithFields(log.Fields{
					"context": "cinema/load",
					"method":  "video/download",
				}),
			)
			return
		}

		log.Println("Time cinema took to trim file", time.Since(n))

		//TODO send file to a fileStore interface. this interface stores the file
		// example fileStores are s3, fs+cache to hold time out e.t.c
		//TODO notify the state of the video with websockets
	}()

	return key, false, nil
}

//getInput returns the input interface that will be used to download the video
//it uses the url to get the input
func (s basicService) getInput(hostname string) input.Interface {
	switch {
	case strings.Contains(hostname, "youtu"):
		return s.youtube
	default:
		return s.url
	}
}

func (s basicService) generateFileKey(d DownloadStruct) (string, bool) {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s-%d-%d", d.URL, d.Start, d.End)))
	key := fmt.Sprintf("%x", hash.Sum(nil))[:10] + "." + d.Type

	//check the fileStore if the file exists
	return key, s.fileStore.Exists(key)
}

func (s basicService) GeneratePath(key string) string {
	return s.fileStore.GeneratePath(key)
}

func (s basicService) Exists(key string) bool {
	return s.fileStore.Exists(key)
}
