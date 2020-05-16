package handler

import (
	"errors"
	"fmt"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"gitlab.com/kayslay/vid_trimmer/pkg/video/service"
	"net/http"
	url2 "net/url"
	"strconv"
	"time"
)

const (
	MaxMinutes = 30
)

func createDownloadStruct(r *http.Request) (service.DownloadStruct, error) {

	q := r.URL.Query()
	url, err := url2.QueryUnescape(q.Get("url"))
	if err != nil {
		return service.DownloadStruct{}, err
	}
	format := q.Get("format")
	if format == "" {
		format = "mp4"
	}

	err2 := validate.Validate(
		&validators.URLIsPresent{Name: "url", Field: url},
		//if check if format is a supported format
		&validators.StringInclusion{
			Name:  "type",
			Field: format,
			List:  []string{"mp4", "mov", "gif", "ogg"},
		})
	if err2.HasAny() {
		return service.DownloadStruct{}, err2
	}

	start, _ := strconv.Atoi(q.Get("start"))
	end, _ := strconv.Atoi(q.Get("end"))

	if end < start {
		return service.DownloadStruct{}, errors.New("end time should be more than start time")
	}

	if end-start > 60*MaxMinutes {
		return service.DownloadStruct{}, errors.New(fmt.Sprintf("max duration of %d min exceeded", MaxMinutes))
	}

	if end-start == 0 {
		return service.DownloadStruct{}, errors.New("duration is 0")
	}

	return service.DownloadStruct{
		URL:   url,
		Start: time.Duration(start) * time.Second,
		End:   time.Duration(end) * time.Second,
		Type:  format,
	}, nil
}
