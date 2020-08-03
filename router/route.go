package router

import (
	"github.com/coda-it/gowebserver/session"
	"github.com/coda-it/gowebserver/store"
	"net/http"
)

// URLOptions - url options type
type URLOptions struct {
	Params map[string]string
}

// ControllerHandler - handler type
type ControllerHandler func(http.ResponseWriter, *http.Request, URLOptions, session.ISessionManager, store.IStore)

// URLRoute - url route type
type URLRoute struct {
	urlRegExp string
	method    string
	handler   ControllerHandler
	params    map[string]int
	protected bool
}
