package main

import (
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/kayslay/vid_trimmer/config"
	"gitlab.com/kayslay/vid_trimmer/pkg/video"
	"net/http"
	"time"
)

func main() {
	godotenv.Load()
	viper.AutomaticEnv()
	port := ":8080"
	if viper.GetString(config.EnvPort) != "" {
		port = ":" + viper.GetString(config.EnvPort)
	}

	svr := http.Server{
		Addr:         port,
		Handler:      http.TimeoutHandler(initRoute(), time.Minute, " server timeout"),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	log.Info("serving http at " + port)
	log.Fatal(svr.ListenAndServe())
}

func initRoute() http.Handler {
	r := chi.NewRouter()
	r.Mount("/download", video.Router())
	return r
}

//"https://media.w3.org/2010/05/sintel/trailer.mp4"
//https://video.twimg.com/ext_tw_video/1261596073178652674/pu/vid/320x568/9H1vwy9y9u9xRXiK.mp4?tag=10
