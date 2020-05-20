package input

import (
	"context"
	yt "github.com/kkdai/youtube"
	"github.com/spf13/viper"
	"gitlab.com/kayslay/vid_trimmer/config"
)

type youtube struct {
	dir      string
	urlInput Interface
}

func NewYoutube(dir string, url Interface) *youtube {
	return &youtube{dir: dir, urlInput: url}
}

func (y youtube) Fetch(ctx context.Context, p string) (string, error) {

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
