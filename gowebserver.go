package gowebserver

import (
    "log"
	"net/http"
    "github.com/oskarszura/trips/gowebserver/router"
    "github.com/oskarszura/trips/gowebserver/session"
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
    session.InitializeSessions()
	staticFileServer := http.FileServer(http.Dir(options.StaticFilesDir))

	http.Handle(options.StaticFilesUrl,
        http.StripPrefix(options.StaticFilesUrl, staticFileServer))
	http.HandleFunc("/", s.Router.Route)

	log.Println("Server listening on port = " + options.Port + " ...")

	err := http.ListenAndServe(options.Port, nil)

	if err != nil {
		log.Fatal("Server failed: ", err)
	}
}
