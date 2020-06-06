package input

import (
	"context"
	"fmt"
	"github.com/bushaHQ/httputil/errors"
	"github.com/dchest/uniuri"
	"github.com/spf13/viper"
	"gitlab.com/kayslay/vid_trimmer/config"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	MbShiftBy = 20
)

type link struct {
	dir string
}

func NewLink(dir string) Interface {
	return &link{dir: dir}
}

func (l link) Fetch(ctx context.Context, p string) (string, error) {

	s := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("could not get video file", 401)
	}

	defer resp.Body.Close()

	cl := resp.Header.Get("Content-Length")
	log.Println("content-length", cl)
	clInt, err := strconv.Atoi(cl)

	//if err != nil {
	//	return "", errors.New("link does not specify file size")
	//}

	if err == nil && clInt > int(getMaxSize()) {
		return "", errors.New(fmt.Sprintf("video is greater than %.4f MB", float64(getMaxSize())/(1<<MbShiftBy)))
	}

	log.Println(resp.Header)

	outputPath := path.Join(l.dir, uniuri.NewLen(10))

	out, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.CopyN(out, resp.Body, getMaxSize())
	if err != nil {
		if err != io.EOF {
			return "", err
		}
	}

	log.Println("GETTER TIME", time.Since(s))

	return outputPath, nil
}

func getMaxSize() int64 {
	return viper.GetInt64(config.EnvFileSize) << MbShiftBy
}
