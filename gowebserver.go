package gowebserver

import (
	"net/http"
    "github.com/oskarszura/gowebserver/router"
    "github.com/oskarszura/gowebserver/session"
	"github.com/oskarszura/gowebserver/utils/logger"
)

type WebServerOptions struct {
    Port            string
    StaticFilesUrl  string
    StaticFilesDir  string
}

type WebServer struct {
	Router		router.UrlRouter
}

func (s *WebServer) Run(options WebServerOptions) {
	logger.Setup("server")

    session.InitializeSessions()
	staticFileServer := http.FileServer(http.Dir(options.StaticFilesDir))

	http.Handle(options.StaticFilesUrl,
        http.StripPrefix(options.StaticFilesUrl, staticFileServer))
	http.HandleFunc("/", s.Router.Route)

	logger.Log("info","Server listening on port = " + options.Port + " ...")

	err := http.ListenAndServe(options.Port, nil)

	if err != nil {
		logger.Log("error","Server failed: ", err)
	}
}
