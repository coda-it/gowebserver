package router

import (
	"github.com/coda-it/gowebserver/session"
	"github.com/coda-it/gowebserver/store"
	"github.com/coda-it/gowebserver/utils/logger"
	"github.com/coda-it/gowebserver/utils/url"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// IRouter - router interface
type IRouter interface {
	Route(w http.ResponseWriter, r *http.Request)
	AddRoute(w http.ResponseWriter, r *http.Request)
}

// Router - router struct
type Router struct {
	sessionManager         session.ISessionManager
	urlRoutes              []URLRoute
	urlRoutesMutex         *sync.RWMutex
	pageNotFoundController ControllerHandler
	store                  store.IStore
	// SessionFallbackURL - when no existing session is detected, router should route here
	SessionFallbackURL string
}

// New - factory for router
func New(sm session.Manager, notFound ControllerHandler, sessionFallbackURL string) Router {
	return Router{
		sessionManager:         sm,
		urlRoutes:              make([]URLRoute, 0),
		urlRoutesMutex:         &sync.RWMutex{},
		pageNotFoundController: notFound,
		store:                  store.New(),
		SessionFallbackURL:     sessionFallbackURL,
	}
}

func (router *Router) findRoute(req *http.Request) URLRoute {
	path := req.URL.Path
	method := req.Method

	router.urlRoutesMutex.RLock()
	urlRoutes := router.urlRoutes
	router.urlRoutesMutex.RUnlock()

	for _, v := range urlRoutes {
		pathRegExp := regexp.MustCompile(v.urlRegExp)

		if pathRegExp.MatchString(path) && (v.method == method || v.method == "ALL") && v.checkerHandler(req) {
			return v
		}
	}
	return URLRoute{
		handler: router.pageNotFoundController,
	}
}

// New - factory for session manager
func (router *Router) New(sm session.ISessionManager) {
	router.sessionManager = sm
}

// Route - routes all incomming requests
func (router *Router) Route(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	route := router.findRoute(r)

	params := make(map[string]string)
	pathItems := strings.Split(urlPath, "/")

	for k, v := range route.params {
		if len(pathItems) > v {
			params[k] = pathItems[v]
		}
	}

	urlOptions := &URLOptions{
		params,
	}

	logger.Log(logger.INFO, "Navigating to url = "+urlPath+" vs route = "+
		route.urlRegExp)

	if route.protected {
		sid, err := session.GetSessionID(r)

		if err != nil || router.sessionManager.IsExist(sid) == false {
			http.Redirect(w, r, router.SessionFallbackURL, http.StatusSeeOther)
			return
		}
	}

	if route.checkerHandler != nil && !route.checkerHandler(r) {
		router.pageNotFoundController(w, r, *urlOptions, router.sessionManager, router.store)
		return
	}

	routeHandler := route.handler
	routeHandler(w, r, *urlOptions, router.sessionManager, router.store)
}

func buildRoute(urlPattern string, method string, protected bool, pathHandler ControllerHandler, checkerHandler CheckerHandler) URLRoute {
	params := make(map[string]int)
	pathRegExp := url.PatternToRegExp(urlPattern)

	urlPathItems := strings.Split(urlPattern, "/")

	for i := 0; i < len(urlPathItems); i++ {
		paramKey := urlPathItems[i]
		isParam, _ := regexp.MatchString(`{[a-zA-Z0-9]*}`, paramKey)

		if isParam {
			strippedParamKey := strings.Replace(strings.Replace(paramKey,
				"{", "", -1), "}", "", -1)
			params[strippedParamKey] = i
		}
	}

	return URLRoute{
		urlRegExp:      pathRegExp,
		method:         method,
		handler:        pathHandler,
		params:         params,
		protected:      protected,
		checkerHandler: checkerHandler,
	}
}

// AddRoute - adds route
func (router *Router) AddRoute(urlPattern string, method string, protected bool, pathHandler ControllerHandler, checkerHandler CheckerHandler) {
	route := buildRoute(urlPattern, method, protected, pathHandler, checkerHandler)

	router.urlRoutesMutex.Lock()
	router.urlRoutes = append(router.urlRoutes, route)
	router.urlRoutesMutex.Unlock()
}

// AddRouteFunc - function used to register a single route during ReplaceRoutes
type AddRouteFunc func(urlPattern string, method string, protected bool, pathHandler ControllerHandler, checkerHandler CheckerHandler)

// ReplaceRoutes - atomically replaces all registered routes with the ones
// added by the register callback; requests keep matching against the old
// routes until the swap happens
func (router *Router) ReplaceRoutes(register func(addRoute AddRouteFunc)) {
	newRoutes := make([]URLRoute, 0)

	register(func(urlPattern string, method string, protected bool, pathHandler ControllerHandler, checkerHandler CheckerHandler) {
		newRoutes = append(newRoutes, buildRoute(urlPattern, method, protected, pathHandler, checkerHandler))
	})

	router.urlRoutesMutex.Lock()
	router.urlRoutes = newRoutes
	router.urlRoutesMutex.Unlock()
}

// AddDataSource - adds data source
func (router *Router) AddDataSource(name string, ds interface{}) {
	router.store.AddDataSource(name, ds)
}
