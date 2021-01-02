package input

import (
	"context"
	"fmt"
	errors2 "github.com/bushaHQ/httputil/errors"
	"github.com/dchest/uniuri"
	yt "github.com/kkdai/youtube/v2"
	log "github.com/sirupsen/logrus"
	"strings"

	"io"
	"net/http"
	"net/url"
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

	u, _ := url.Parse(p)
	videoID := u.Query().Get("v")
	if videoID == "" {
		videoID = strings.Replace(u.Path, "/", "", -1)
	}

	client := yt.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return "", err
	}

	resp, err := client.GetStream(video, video.Formats.FindByQuality("medium"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	outputPath := path.Join(y.dir, uniuri.NewLen(10))

	out, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := copy(out, resp.Body); err != nil {
		if err != io.EOF {
			return "", err
		}
	}

	return outputPath, nil
}
