package handler

import (
	"fmt"
	"github.com/bushaHQ/httputil/errors"
	"github.com/bushaHQ/httputil/render"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"gitlab.com/kayslay/vid_trimmer/config"
	"gitlab.com/kayslay/vid_trimmer/internal/filestore"
	"gitlab.com/kayslay/vid_trimmer/pkg/video/service"
	"io"
	"mime"
	"net/http"
	"path"
)

type Video struct {
	svc service.Service
}

func NewVideo(svc service.Service) *Video {
	return &Video{svc: svc}
}

func (h Video) GenerateDownloadLink(w http.ResponseWriter, r *http.Request) {
	ds, err := createDownloadStruct(r)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	key, state, err := h.svc.Download(r.Context(), ds)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	render.Render(w, r, map[string]interface{}{
		"link":  viper.GetString(config.EnvHost) + "/download/" + key,
		"state": state,
	})

}

func (h Video) Download(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "key")

	state := h.svc.FileState(key)

	if state == filestore.StateNull {
		render.Render(w, r, errors.New("file not does not exists", http.StatusNotFound))
		return
	}

	if state == filestore.StatePending {
		render.Render(w, r, errors.New("wait, I'm still trimming the file", http.StatusTeapot))
		return
	}

	file, fileStat, err := h.svc.Get(key)
	if err != nil {
		render.Render(w, r, errors.New("could not get file at the moment", 404))
		return
	}

	defer file.Close()

	w.Header().Set("Content-Type", mime.TypeByExtension("."+path.Ext(key)))
	w.Header().Set("Content-Disposition", "inline; filename="+key)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size))

	io.Copy(w, file)
}
