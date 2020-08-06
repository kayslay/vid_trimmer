package input

import (
	"context"
	"fmt"
	errors2 "github.com/bushaHQ/httputil/errors"
	"github.com/dchest/uniuri"
	log "github.com/sirupsen/logrus"
	yt "github.com/wader/goutubedl"
	"io"
	"net/http"
	"os"
	"path"
)

type youtube struct {
	dir      string
	urlInput Interface
}

func NewYoutube(dir string, url Interface) *youtube {
	return &youtube{dir: dir, urlInput: url}
}

func (y youtube) Fetch(ctx context.Context, p string) (str string, err error) {

	defer func() {
		if rvr := recover(); rvr != nil {
			err = errors2.CoverErr(
				fmt.Errorf("%v", rvr),
				errors2.New("Can not trim this video file. working on a fix", http.StatusServiceUnavailable),
				log.WithFields(log.Fields{
					"context": "defer",
					"method":  "youtube/fetch",
				}),
			)
		}
	}()

	result, err := yt.New(ctx, p, yt.Options{})
	if err != nil {
		return "", err
	}

	log.Println(string(result.RawJSON))
	downloadResult, err := result.Download(ctx, "best[height<=480]")
	if err != nil {
		return "", err
	}

	defer downloadResult.Close()

	outputPath := path.Join(y.dir, uniuri.NewLen(10))

	out, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.CopyN(out, downloadResult, getMaxSize())
	if err != nil {
		if err != io.EOF {
			return "", err
		}
	}

	return outputPath, nil
}
