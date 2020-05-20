package handler

import (
	"bytes"
	"fmt"
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

	fileName, err := h.svc.Download(r.Context(), buf, ds)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension("."+ds.Type))
	w.Header().Set("Content-Disposition", "inline; filename="+fileName)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	buf.WriteTo(w)
}
