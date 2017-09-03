package router

import "net/http"

type UrlOptions struct {
    Params map[string]string
}

type ControllerHandler func(http.ResponseWriter, *http.Request, UrlOptions)

type UrlRoute struct {
    urlRegExp string
    method    string
    handler	  ControllerHandler
    params	  map[string]int
}
