package gowebserver

import (
	"regexp"
	"net/http"
)

type ControllerHandler func(http.ResponseWriter, *http.Request)

type UrlRoute struct {
	urlRegExp string
	method    string
	handler	  ControllerHandler
}

type UrlRouter struct {
	urlRoutes 		[]UrlRoute
	pageNotFoundController	ControllerHandler
}

func (router *UrlRouter) findRoute (path string) UrlRoute {
	for _, v := range router.urlRoutes {
		pathRegExp := regexp.MustCompile(v.urlRegExp)

		if (pathRegExp.MatchString(path)) {
			return v
		}
	}
	return UrlRoute{
		handler: router.pageNotFoundController,
	}
}

func (router *UrlRouter) AddRoute(pathRegExp string, pathHandler ControllerHandler) {
	router.urlRoutes = append(router.urlRoutes, UrlRoute{
		urlRegExp: pathRegExp,
		handler: pathHandler,
	})
}

func (router *UrlRouter) AddPageNotFoundRoute(pathHandler ControllerHandler) {
	router.pageNotFoundController = pathHandler
}

func (router *UrlRouter) route(w http.ResponseWriter, r *http.Request)  {
	urlPath := r.URL.Path

	routeHandler := router.findRoute(urlPath).handler
	routeHandler(w, r)
}
