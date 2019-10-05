package router

import (
	"github.com/coda-it/gowebserver/session"
	"github.com/coda-it/gowebserver/store"
	"net/http"
)

type UrlOptions struct {
	Params map[string]string
}

type ControllerHandler func(http.ResponseWriter, *http.Request, UrlOptions, session.ISessionManager, store.IStore)

type UrlRoute struct {
	urlRegExp string
	method    string
	handler   ControllerHandler
	params    map[string]int
}
