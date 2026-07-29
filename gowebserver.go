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

	// http.Server timeouts. A zero value applies a safe default; a negative
	// value disables the timeout. Disable WriteTimeout (set it negative) for
	// long-lived streaming responses such as SSE, where a fixed write deadline
	// would otherwise cut the stream off.
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

const (
	defaultReadHeaderTimeout = 10 * time.Second
	defaultReadTimeout       = 30 * time.Second
	defaultWriteTimeout      = 60 * time.Second
	defaultIdleTimeout       = 120 * time.Second
)

// resolveTimeout maps a configured duration to an http.Server timeout: zero
// selects the default, negative disables the timeout (0 for net/http).
func resolveTimeout(configured, def time.Duration) time.Duration {
	switch {
	case configured < 0:
		return 0
	case configured == 0:
		return def
	default:
		return configured
	}
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

// cacheControlResponseWriter applies the cache policy only once the status is
// known, so error responses (404, 5xx, ...) are not marked publicly cacheable.
type cacheControlResponseWriter struct {
	http.ResponseWriter
	value       string
	wroteHeader bool
}

func (w *cacheControlResponseWriter) WriteHeader(status int) {
	if !w.wroteHeader {
		w.wroteHeader = true
		if status < 400 {
			w.Header().Set("Cache-Control", w.value)
		}
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *cacheControlResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func withStaticCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&cacheControlResponseWriter{ResponseWriter: w, value: staticCacheControl}, r)
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
		ReadHeaderTimeout: resolveTimeout(s.Options.ReadHeaderTimeout, defaultReadHeaderTimeout),
		ReadTimeout:       resolveTimeout(s.Options.ReadTimeout, defaultReadTimeout),
		WriteTimeout:      resolveTimeout(s.Options.WriteTimeout, defaultWriteTimeout),
		IdleTimeout:       resolveTimeout(s.Options.IdleTimeout, defaultIdleTimeout),
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
