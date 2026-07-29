package middleware

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAcceptsGzip(t *testing.T) {
	cases := map[string]bool{
		"":                       false,
		"identity":               false,
		"gzip":                   true,
		"gzip, deflate":          true,
		"deflate, gzip;q=0.8":    true,
		"gzip;q=0":               false,
		"gzip;Q=0":               false, // parameter name is case-insensitive
		"gzip; q=0":              false,
		"gzip;q=0.0":             false,
		"*":                      true,
		"*;q=0":                  false,
		"deflate, *;q=0":         false,
		"gzip;q=0, *":            false, // explicit gzip refusal wins over wildcard
		"br;q=1.0, gzip;q=0.001": true,
	}
	for header, want := range cases {
		if got := acceptsGzip(header); got != want {
			t.Errorf("acceptsGzip(%q) = %v, want %v", header, got, want)
		}
	}
}

func gzipDecode(t *testing.T, b []byte) string {
	t.Helper()
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("gzip.NewReader: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("gzip read: %v", err)
	}
	return string(out)
}

func TestGzipCompressesWhenAccepted(t *testing.T) {
	body := strings.Repeat("hello world ", 100)
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, body)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	res := rec.Result()
	if res.Header.Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected gzip encoding, got %q", res.Header.Get("Content-Encoding"))
	}
	if !strings.Contains(res.Header.Get("Vary"), "Accept-Encoding") {
		t.Errorf("missing Vary: Accept-Encoding")
	}
	if got := gzipDecode(t, rec.Body.Bytes()); got != body {
		t.Errorf("decoded body mismatch")
	}
}

func TestGzipSkippedWhenRefused(t *testing.T) {
	body := "plain body"
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, body)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip;q=0")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	res := rec.Result()
	if enc := res.Header.Get("Content-Encoding"); enc != "" {
		t.Fatalf("expected no encoding, got %q", enc)
	}
	if rec.Body.String() != body {
		t.Errorf("body = %q, want %q", rec.Body.String(), body)
	}
	// Vary must be present even on the identity path so shared caches do not
	// reuse this entry for gzip-capable clients.
	if !strings.Contains(res.Header.Get("Vary"), "Accept-Encoding") {
		t.Errorf("missing Vary: Accept-Encoding on bypass path")
	}
}

func TestGzipSkipsUncompressibleContentType(t *testing.T) {
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte("\x89PNG\r\n\x1a\n rest of the bytes"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if enc := rec.Result().Header.Get("Content-Encoding"); enc != "" {
		t.Fatalf("image should not be gzipped, got encoding %q", enc)
	}
}

func TestGzipSniffsContentTypeAfterWriteHeader(t *testing.T) {
	// Controller sets a status before writing and never sets Content-Type;
	// the body should still be sniffed just like an unwrapped response.
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "<!DOCTYPE html><html></html>")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if ct := rec.Result().Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected sniffed text/html Content-Type, got %q", ct)
	}
}

func TestGzipEmptyBodyNotEncoded(t *testing.T) {
	// A handler that writes no body must not advertise gzip on an empty payload.
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", res.StatusCode)
	}
	if enc := res.Header.Get("Content-Encoding"); enc != "" {
		t.Errorf("empty body should not be gzip-encoded, got %q", enc)
	}
}

func TestIsCompressible(t *testing.T) {
	cases := map[string]bool{
		"text/html":                    true,
		"text/html; charset=utf-8":     true,
		"application/json":             true,
		"image/png":                    false,
		"image/svg+xml":                true,
		"image/svg+xml; charset=utf-8": true, // parameterised SVG still compresses
		"video/mp4":                    false,
		"font/woff2":                   false,
		"":                             true, // unset: default to compressing
	}
	for ct, want := range cases {
		if got := isCompressible(ct); got != want {
			t.Errorf("isCompressible(%q) = %v, want %v", ct, got, want)
		}
	}
}

func TestGzipPreservesHandlerVary(t *testing.T) {
	// A handler that sets Vary itself (via Set, which replaces) must not lose
	// Accept-Encoding from the committed response.
	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, strings.Repeat("x", 200))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	// Vary may be sent as multiple header lines; caches combine them, so check
	// the full set rather than just the first line.
	vary := strings.Join(rec.Result().Header["Vary"], ", ")
	if !strings.Contains(vary, "Origin") || !strings.Contains(vary, "Accept-Encoding") {
		t.Errorf("Vary = %q, want both Origin and Accept-Encoding", vary)
	}
}

// fullRecorder augments httptest.ResponseRecorder with the optional interfaces
// a real net/http writer exposes, so we can assert the wrapper forwards them.
type fullRecorder struct {
	*httptest.ResponseRecorder
	flushed    bool
	hijacked   bool
	closeNotCh chan bool
}

func (r *fullRecorder) Flush() { r.flushed = true }

func (r *fullRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	r.hijacked = true
	return nil, nil, nil
}

func (r *fullRecorder) CloseNotify() <-chan bool { return r.closeNotCh }

func TestGzipForwardsOptionalInterfaces(t *testing.T) {
	var sawFlusher, sawHijacker, sawCloseNotifier bool
	rec := &fullRecorder{ResponseRecorder: httptest.NewRecorder(), closeNotCh: make(chan bool)}

	h := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if f, ok := w.(http.Flusher); ok {
			sawFlusher = true
			io.WriteString(w, "data")
			f.Flush()
		}
		if hj, ok := w.(http.Hijacker); ok {
			sawHijacker = true
			hj.Hijack()
		}
		if cn, ok := w.(http.CloseNotifier); ok {
			// Forwarded channel must be the underlying one.
			if cn.CloseNotify() == rec.closeNotCh {
				sawCloseNotifier = true
			}
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	h.ServeHTTP(rec, req)

	if !sawFlusher {
		t.Error("handler could not assert http.Flusher on the wrapped writer")
	}
	if !sawHijacker {
		t.Error("handler could not assert http.Hijacker on the wrapped writer")
	}
	if !sawCloseNotifier {
		t.Error("http.CloseNotifier was not forwarded to the underlying writer")
	}
	if !rec.flushed {
		t.Error("Flush was not forwarded to the underlying writer")
	}
	if !rec.hijacked {
		t.Error("Hijack was not forwarded to the underlying writer")
	}
}
