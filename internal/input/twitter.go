package input

import (
	"context"
	"github.com/bushaHQ/httputil/errors"
	tw "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
	config2 "gitlab.com/kayslay/vid_trimmer/config"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type twitter struct {
	dir      string
	urlInput Interface
	cl       *tw.Client
}

func NewTwitter(dir string, urlInput Interface) *twitter {
	config := oauth1.NewConfig(viper.GetString(config2.EnvTwitterConsumerKey), viper.GetString(config2.EnvTwitterConsumerSecret))
	token := oauth1.NewToken(viper.GetString(config2.EnvTwitterAccessKey), viper.GetString(config2.EnvTwitterAccessSecret))
	client := config.Client(oauth1.NoContext, token)
	return &twitter{dir: dir, urlInput: urlInput, cl: tw.NewClient(client)}
}

func (t twitter) Fetch(ctx context.Context, p string) (string, error) {
	u, err := url.Parse(p)
	if err != nil {
		return "", err
	}
	id, _ := strconv.ParseInt(path.Base(u.Path), 10, 64)
	showEntities := true
	tweet, _, err := t.cl.Statuses.Show(id, &tw.StatusShowParams{
		ID:              id,
		IncludeEntities: &showEntities,
		TweetMode:       "extended",
	})
	if err != nil {
		return "", err
	}

	if len(tweet.Entities.Media) == 0 {
		return "", errors.New("tweet has no media", http.StatusBadRequest)
	}

	media := tweet.ExtendedEntities.Media[0]
	var videoLink string
	for _, variant := range media.VideoInfo.Variants {
		if variant.ContentType == "video/mp4" {
			videoLink = variant.URL
		}
	}

	return t.urlInput.Fetch(ctx, videoLink)
}
