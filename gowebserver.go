package gowebserver

import (
    "net/http"
    "github.com/coda-it/gowebserver/router"
    "github.com/coda-it/gowebserver/session"
    "github.com/coda-it/gowebserver/utils/logger"
)

type WebServerOptions struct {
    Port            string
    StaticFilesUrl  string
    StaticFilesDir  string
}

type WebServer struct {
	Router			router.Router
	Options 		WebServerOptions
}

func New(options WebServerOptions, notFound router.ControllerHandler) *WebServer {
	sm := session.New()

	return &WebServer{
		router.New(sm, notFound),
		options,
	}
}

func (s *WebServer) Run() bool {
	logger.Init("server")

	staticFileServer := http.FileServer(http.Dir(s.Options.StaticFilesDir))

	http.Handle(s.Options.StaticFilesUrl, http.StripPrefix(s.Options.StaticFilesUrl, staticFileServer))
	http.HandleFunc("/", s.Router.Route)

	logger.Log(logger.INFO,"Server listening on port = " + s.Options.Port+ " ...")

	err := http.ListenAndServe(s.Options.Port, nil)

	if err != nil {
		logger.Log(logger.INFO,"Running server failed: ", err)
		return false
	}

	return true
}

func (s *WebServer) AddDataSource(name string, ds interface{}) {
	s.Router.AddDataSource(name, ds)
}