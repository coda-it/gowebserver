package gowebserver

import (
	"github.com/coda-it/gowebserver/router"
	"github.com/coda-it/gowebserver/session"
	"github.com/coda-it/gowebserver/utils/logger"
	"net/http"
)

// WebServerOptions - web server options
type WebServerOptions struct {
	Port           string
	StaticFilesURL string
	StaticFilesDir string
}

// WebServer - web server
type WebServer struct {
	Router  router.Router
	Options WebServerOptions
}

// New - factory for WebServer entity
func New(options WebServerOptions, notFound router.ControllerHandler, sessionFallbackURL string) *WebServer {
	sm := session.New()

	return &WebServer{
		router.New(sm, notFound, sessionFallbackURL),
		options,
	}
}

// Run - starts WebServer
func (s *WebServer) Run() bool {
	logger.Init("server")

	staticFileServer := http.FileServer(http.Dir(s.Options.StaticFilesDir))

	http.Handle(s.Options.StaticFilesURL, http.StripPrefix(s.Options.StaticFilesURL, staticFileServer))
	http.HandleFunc("/", s.Router.Route)

	logger.Log(logger.INFO, "Server listening on port = "+s.Options.Port+" ...")

	err := http.ListenAndServe(s.Options.Port, nil)

	if err != nil {
		logger.Log(logger.INFO, "Running server failed: ", err)
		return false
	}

	return true
}

// AddDataSource - adds data source to WebServer
func (s *WebServer) AddDataSource(name string, ds interface{}) {
	s.Router.AddDataSource(name, ds)
}
