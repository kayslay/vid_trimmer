package handler

import (
	"bytes"
	"github.com/bushaHQ/httputil/render"
	"gitlab.com/kayslay/vid_trimmer/pkg/video/service"
	"mime"
	"net/http"
)

type Video struct {
	svc service.Service
}

func NewVideo(svc service.Service) *Video {
	return &Video{svc: svc}
}

func (h Video) Download(w http.ResponseWriter, r *http.Request) {
	ds, err := createDownloadStruct(r)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	buf := bytes.NewBuffer([]byte{})

	err = h.svc.Download(buf, ds)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension("."+ds.Type))
	buf.WriteTo(w)
}
