package input

import (
	"context"
	"fmt"
	errors2 "github.com/bushaHQ/httputil/errors"
	yt "github.com/kkdai/youtube"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/kayslay/vid_trimmer/config"
	"net/http"
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

	dl := yt.NewYoutube(viper.GetString(config.Env) != "production")
	if err := dl.DecodeURL(p); err != nil {
		return "", err
	}

	//outputPath := path.Join(y.dir, uniuri.NewLen(10))

	var vUrl string
	for _, stream := range dl.StreamList {

		if stream["quality"] == "medium" {
			vUrl = stream["url"]
		}
	}

	return y.urlInput.Fetch(ctx, vUrl)

	//if err := dl.StartDownloadWithQuality(outputPath, "medium"); err != nil {
	//	return "", err
	//}
	//
	//return outputPath, nil
}
