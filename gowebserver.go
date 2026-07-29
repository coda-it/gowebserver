package gowebserver

import (
	"github.com/coda-it/gowebserver/middleware"
	"github.com/coda-it/gowebserver/router"
	"github.com/coda-it/gowebserver/session"
	"github.com/coda-it/gowebserver/utils/logger"
	"net/http"
	"time"
)

// WebServerOptions - web server options
type WebServerOptions struct {
	Port           string
	StaticFilesURL string
	StaticFilesDir string
	Cert           string
	CertPrvKey     string
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

// staticCacheControl - assets are not content-hashed, so allow short-lived
// caching and rely on ETag/Last-Modified revalidation from http.FileServer
const staticCacheControl = "public, max-age=3600, must-revalidate"

func withStaticCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", staticCacheControl)
		next.ServeHTTP(w, r)
	})
}

// Run - starts WebServer
func (s *WebServer) Run() bool {
	logger.Init("server")

	staticFileServer := http.FileServer(http.Dir(s.Options.StaticFilesDir))

	http.Handle(s.Options.StaticFilesURL, middleware.Gzip(withStaticCacheControl(
		http.StripPrefix(s.Options.StaticFilesURL, staticFileServer))))
	http.Handle("/", middleware.Gzip(http.HandlerFunc(s.Router.Route)))

	logger.Log(logger.INFO, "Server listening on port = "+s.Options.Port+" ...")

	server := &http.Server{
		Addr:              s.Options.Port,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	var err error
	if s.Options.Cert != "" && s.Options.CertPrvKey != "" {
		err = server.ListenAndServeTLS(s.Options.Cert, s.Options.CertPrvKey)
	} else {
		err = server.ListenAndServe()
	}

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
