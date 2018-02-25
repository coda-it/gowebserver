package router

import (
    "net/http"
    "github.com/oskarszura/gowebserver/session"
)

type UrlOptions struct {
    Params map[string]string
}

type ControllerHandler func(http.ResponseWriter, *http.Request, UrlOptions, session.ISessionManager)

type UrlRoute struct {
    urlRegExp       string
    method          string
    handler	        ControllerHandler
    params	        map[string]int
}
