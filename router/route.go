package router

import (
    "net/http"
    "github.com/coda-it/gowebserver/session"
    "github.com/coda-it/gowebserver/store"
)

type UrlOptions struct {
    Params map[string]string
}

type ControllerHandler func(http.ResponseWriter, *http.Request, UrlOptions, session.ISessionManager, store store.IStore)

type UrlRoute struct {
    urlRegExp       string
    method          string
    handler	        ControllerHandler
    params	        map[string]int
}
