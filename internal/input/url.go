package input

import (
	"context"
	"fmt"
	"github.com/bushaHQ/httputil/errors"
	"github.com/dchest/uniuri"
	"github.com/spf13/viper"
	"gitlab.com/kayslay/vid_trimmer/config"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("could not get video file", 401)
	}

	defer resp.Body.Close()

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))

	if err == nil && contentLength > int(getMaxSize()) {
		return "", errors.New(fmt.Sprintf("video is greater than %.4f MB", float64(getMaxSize())/(1<<MbShiftBy)))
	}

	//set output path for video
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

	return outputPath, nil
}

func getMaxSize() int64 {
	size := viper.GetInt64(config.EnvFileSize)
	if size == 0 {
		return 0
	}
	return size << MbShiftBy
}

//copy copies the content of the reader to the writer. copies all or part depending config.EnvFileSize
func copy(out io.Writer, r io.Reader) error {
	var err error
	maxSize := getMaxSize()

	if maxSize > 0 {
		_, err = io.CopyN(out, r, maxSize)
	} else {
		_, err = io.Copy(out, r)
	}

	return err
}
