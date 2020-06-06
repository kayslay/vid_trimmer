package handler

import (
	"fmt"
	"github.com/bushaHQ/httputil/errors"
	"github.com/bushaHQ/httputil/render"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"gitlab.com/kayslay/vid_trimmer/config"
	"gitlab.com/kayslay/vid_trimmer/pkg/video/service"
	"io"
	"mime"
	"net/http"
	"os"
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

	key, exist, err := h.svc.Download(r.Context(), ds)
	if err != nil {
		render.Render(w, r, err)
		return
	}

	render.Render(w, r, map[string]interface{}{
		"link":   viper.GetString(config.EnvHost) + "/download/" + key,
		"exists": exist,
	})

}

func (h Video) Download(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "key")

	if !h.svc.Exists(key) {
		render.Render(w, r, errors.New("file not does not exists", 404))
		return
	}

	file, err := os.Open(h.svc.GeneratePath(key))
	if err != nil {
		render.Render(w, r, errors.New("could not get file at the moment", 404))
		return
	}

	defer file.Close()
	fileInfo, _ := file.Stat()

	w.Header().Set("Content-Type", mime.TypeByExtension("."+path.Ext(key)))
	w.Header().Set("Content-Disposition", "inline; filename="+key)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	io.Copy(w, file)
}
