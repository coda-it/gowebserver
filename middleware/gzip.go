package middleware

import (
	"bufio"
	"compress/gzip"
	"errors"
	"mime"
	"net"
	"net/http"
	"strconv"
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
	// Strip any parameters (e.g. "; charset=utf-8") and normalise case so a
	// parameterised type such as "image/svg+xml; charset=utf-8" is classified
	// by its media type rather than by raw-string matching.
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		// Unknown or unset type: default to compressing (text-like payloads).
		return true
	}
	if mediaType == "image/svg+xml" {
		return true
	}
	for _, prefix := range uncompressibleContentTypes {
		if strings.HasPrefix(mediaType, prefix) {
			return false
		}
	}
	return true
}

// ensureVary adds value to the Vary header unless it (or "*") is already
// present, preserving any Vary fields the handler set itself. It reads the
// header map directly rather than using Header.Values, which was added in Go
// 1.14, for the module's Go 1.12 baseline.
func ensureVary(h http.Header, value string) {
	for _, existing := range h["Vary"] {
		for _, field := range strings.Split(existing, ",") {
			field = strings.TrimSpace(field)
			if field == "*" || strings.EqualFold(field, value) {
				return
			}
		}
	}
	h.Add("Vary", value)
}

// acceptsGzip reports whether the Accept-Encoding header actually allows gzip.
// It parses the codings and their quality values so explicit refusals such as
// "gzip;q=0" (and wildcard entries like "*;q=0") are honoured, rather than
// matching gzip anywhere in the raw header string.
func acceptsGzip(header string) bool {
	if header == "" {
		return false
	}

	wildcardSeen := false
	wildcardOK := false

	for _, part := range strings.Split(header, ",") {
		tokens := strings.Split(part, ";")
		coding := strings.ToLower(strings.TrimSpace(tokens[0]))
		if coding != "gzip" && coding != "*" {
			continue
		}

		q := 1.0
		for _, t := range tokens[1:] {
			// Parameter names are case-insensitive per the HTTP grammar, so a
			// valid "gzip;Q=0" must be recognised as a refusal.
			t = strings.ToLower(strings.TrimSpace(t))
			if strings.HasPrefix(t, "q=") {
				if v, err := strconv.ParseFloat(strings.TrimSpace(t[2:]), 64); err == nil {
					q = v
				}
			}
		}

		// An explicit gzip entry wins over any wildcard.
		if coding == "gzip" {
			return q > 0
		}

		wildcardSeen = true
		wildcardOK = q > 0
	}

	if wildcardSeen {
		return wildcardOK
	}
	return false
}

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter    *gzip.Writer
	status        int
	allowCompress bool // client advertised gzip support
	wroteHeader   bool // caller has chosen a status
	committed     bool // headers have been flushed to the underlying writer
	compress      bool // final decision: this response is being gzipped
}

// WriteHeader records the status but defers writing headers to the underlying
// writer until the first Write, so the body can still be content-sniffed.
func (w *gzipResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.status = status
}

// commit decides whether to compress and flushes headers exactly once. sniff is
// the first chunk of the body (may be nil) used for Content-Type detection.
func (w *gzipResponseWriter) commit(sniff []byte) {
	if w.committed {
		return
	}
	w.committed = true

	status := w.status
	if !w.wroteHeader {
		status = http.StatusOK
	}

	// Representation depends on Accept-Encoding regardless of whether this
	// particular response ends up compressed, so declare it at commit time.
	// Doing it here (not before the handler runs) survives a handler that
	// overwrites Vary via Header().Set.
	ensureVary(w.Header(), "Accept-Encoding")

	if len(sniff) > 0 && w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(sniff))
	}

	// Only compress when the client accepts gzip and there is a body to
	// compress: an empty response must not advertise Content-Encoding: gzip
	// with a zero-length payload.
	w.compress = w.allowCompress &&
		len(sniff) > 0 &&
		status != http.StatusNoContent &&
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
	if !w.committed {
		w.commit(b)
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

// Flush forwards to the underlying writer so streaming responses (e.g. SSE)
// keep working, flushing any buffered compressed data first.
func (w *gzipResponseWriter) Flush() {
	if !w.committed {
		w.commit(nil)
	}
	if w.gzipWriter != nil {
		w.gzipWriter.Flush()
	}
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack forwards to the underlying writer so connection upgrades (e.g.
// WebSocket) keep working when the server supports them.
func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("gzip: underlying ResponseWriter does not implement http.Hijacker")
}

// CloseNotify forwards the underlying writer's CloseNotifier when available.
// http.CloseNotifier is deprecated in favour of Request.Context, but is
// preserved here for legacy SSE/long-poll handlers targeting the module's Go
// 1.12 baseline. Returns a nil channel (which never fires) when the underlying
// writer does not implement it.
func (w *gzipResponseWriter) CloseNotify() <-chan bool {
	if cn, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}

func (w *gzipResponseWriter) close() {
	// Flush headers for responses that never wrote a body (e.g. 204/304).
	if !w.committed {
		w.commit(nil)
	}
	if w.gzipWriter == nil {
		return
	}
	w.gzipWriter.Close()
	gzipWriterPool.Put(w.gzipWriter)
	w.gzipWriter = nil
}

// Gzip - wraps a handler with gzip compression for clients that accept it.
// Both gzip-capable and identity responses flow through the same wrapper so
// optional interfaces (Flusher, Hijacker, CloseNotifier) and the
// Vary: Accept-Encoding header are handled consistently on every path.
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gzw := &gzipResponseWriter{
			ResponseWriter: w,
			allowCompress:  acceptsGzip(r.Header.Get("Accept-Encoding")),
		}
		defer gzw.close()

		next.ServeHTTP(gzw, r)
	})
}
