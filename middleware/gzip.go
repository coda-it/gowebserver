package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

var uncompressibleContentTypes = []string{
	"image/",
	"video/",
	"audio/",
	"font/",
	"application/font-woff",
	"application/vnd.ms-fontobject",
	"application/zip",
	"application/gzip",
	"application/x-gzip",
}

func isCompressible(contentType string) bool {
	if contentType == "image/svg+xml" {
		return true
	}
	for _, prefix := range uncompressibleContentTypes {
		if strings.HasPrefix(contentType, prefix) {
			return false
		}
	}
	return true
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter  *gzip.Writer
	wroteHeader bool
	compress    bool
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true

	w.compress = status != http.StatusNoContent &&
		status != http.StatusNotModified &&
		status != http.StatusPartialContent &&
		w.Header().Get("Content-Encoding") == "" &&
		isCompressible(w.Header().Get("Content-Type"))

	if w.compress {
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Encoding", "gzip")
	}

	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(b))
		}
		w.WriteHeader(http.StatusOK)
	}

	if !w.compress {
		return w.ResponseWriter.Write(b)
	}

	if w.gzipWriter == nil {
		w.gzipWriter = gzipWriterPool.Get().(*gzip.Writer)
		w.gzipWriter.Reset(w.ResponseWriter)
	}

	return w.gzipWriter.Write(b)
}

func (w *gzipResponseWriter) close() {
	if w.gzipWriter == nil {
		return
	}
	w.gzipWriter.Close()
	gzipWriterPool.Put(w.gzipWriter)
	w.gzipWriter = nil
}

// Gzip - wraps a handler with gzip compression for clients that accept it
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Accept-Encoding")

		gzw := &gzipResponseWriter{ResponseWriter: w}
		defer gzw.close()

		next.ServeHTTP(gzw, r)
	})
}
