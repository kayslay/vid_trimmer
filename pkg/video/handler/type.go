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
	"strings"
	"time"
	"unicode"
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

	start, _ := genSecFromDuration(q.Get("start"))
	end, _ := genSecFromDuration(q.Get("end"))

	if end <= start {
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

//genSecFromDuration generates the seconds from the duration passed. if the string does not end with a digit
// return the value
func genSecFromDuration(t string) (int, error) {
	t = strings.TrimSpace(t)

	if len(t) == 0 {
		return 0, nil
	}
	//if the last value is a digit use the value passed
	if unicode.IsNumber(rune(t[len(t)-1])) {
		return strconv.Atoi(t)
	}

	d, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}

	return int(d.Seconds()), nil
}
