package gowebserver

import (
	"regexp"
	"net/http"
	"time"
)

type ControllerHandler func(http.ResponseWriter, *http.Request)

type UrlRoute struct {
	urlRegExp string
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

	_, err := r.Cookie("sid")

	if(err != nil) {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		sid := SidGenerator()

		cookie := http.Cookie {
			Name: "sid",
			Value: sid,
			Expires: expiration }

		http.SetCookie(w, &cookie)
	}

	routeHandler := router.findRoute(urlPath).handler
	routeHandler(w, r)
}
