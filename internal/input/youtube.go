package input

import (
	"github.com/dchest/uniuri"
	yt "github.com/kkdai/youtube"
	"github.com/spf13/viper"
	"gitlab.com/kayslay/vid_trimmer/config"
	"path"
)

type youtube struct {
	dir string
}

func NewYoutube(dir string) *youtube {
	return &youtube{dir: dir}
}

func (y youtube) Fetch(p string) (string, error) {

	dl := yt.NewYoutube(viper.GetString(config.Env) != "production")
	if err := dl.DecodeURL(p); err != nil {
		return "", err
	}

	outputPath := path.Join(y.dir, uniuri.NewLen(10))

	if err := dl.StartDownloadWithQuality(outputPath, "medium"); err != nil {
		return "", err
	}

	return outputPath, nil
}
