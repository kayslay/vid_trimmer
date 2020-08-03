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
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service interface {
	Download(ctx context.Context, d DownloadStruct) (string, string, error)
	Get(key string) (io.ReadCloser, filestore.FileStat, error)
	FileState(key string) string
}

type basicService struct {
	url       input.Interface
	youtube   input.Interface
	fileStore filestore.FileStorer
	twitter   input.Interface
}

func NewBasicService(url, youtube, twitter input.Interface, storer filestore.FileStorer) Service {
	return &basicService{url: url, youtube: youtube, twitter: twitter, fileStore: storer}
}

func (s basicService) Download(ctx context.Context, d DownloadStruct) (string, string, error) {
	key, state := s.generateFileKey(d)

	//if the file is not null it means it pending or done
	//return the key and state
	if state != filestore.StateNull {
		return key, state, nil
	}

	//generate a temp file
	tempOutputPath := filestore.NewTempFile().GeneratePath(key)

	// get the url struct
	u, err := url.Parse(d.URL)
	if err != nil {
		return "", filestore.StateNull, err
	}

	n := time.Now()

	pathUrl, err := s.getInput(u.Hostname()).Fetch(ctx, d.URL)

	if err != nil {
		return "", filestore.StateNull, errors.CoverErr(
			err,
			errors.New("could not create file at the moment", http.StatusServiceUnavailable),
			log.WithFields(log.Fields{
				"context": "fetch",
				"method":  "video/download",
			}),
		)
	}
	log.Println("Time taken to get input", time.Since(n))

	//trim in the background
	go func() {
		s.fileStore.ChangeState(key, filestore.StatePending)
		defer input.Remove(pathUrl)
		defer func() {
			//handle panic
			if err := recover(); err != nil {
				s.fileStore.ChangeState(key, filestore.StateNull)

				log.WithFields(log.Fields{
					"context": "cinema/load",
					"method":  "video/download",
				})
			}
		}()

		n := time.Now()
		v, err := cinema.Load(pathUrl)
		if err != nil {
			s.fileStore.ChangeState(key, filestore.StateNull)

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
		err = v.Render(tempOutputPath)

		if err != nil {
			s.fileStore.ChangeState(key, filestore.StateNull)

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

		//call a file store to store the video s3, tempDir e.t.c
		s.fileStore.Write(key, tempOutputPath)

		//TODO notify the state of the video with websockets
	}()

	return key, filestore.StateNull, nil
}

//getInput returns the input interface that will be used to download the video
//it uses the url to get the input
func (s basicService) getInput(hostname string) input.Interface {
	switch {
	case strings.Contains(hostname, "youtu"):
		return s.youtube
	case strings.Contains(hostname, "twitter"):
		return s.twitter
	default:
		return s.url
	}
}

func (s basicService) generateFileKey(d DownloadStruct) (string, string) {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s-%s-%d-%d", d.URL, d.Type, d.Start, d.End)))
	key := fmt.Sprintf("%x", hash.Sum(nil))[:10] + "." + d.Type

	//check the fileStore if the file exists
	return key, s.fileStore.FileState(key)
}

func (s basicService) Get(key string) (io.ReadCloser, filestore.FileStat, error) {
	return s.fileStore.Get(key)
}

func (s basicService) FileState(key string) string {
	return s.fileStore.FileState(key)
}
