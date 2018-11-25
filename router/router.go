package router

import (
	"regexp"
	"strings"
	"net/http"
	"github.com/coda-it/gowebserver/utils/url"
	"github.com/coda-it/gowebserver/utils/logger"
	"github.com/coda-it/gowebserver/session"
)

type IRouter interface {
	Route(w http.ResponseWriter, r *http.Request)
	AddRoute(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	sessionManager			session.ISessionManager
    urlRoutes 				[]UrlRoute
    pageNotFoundController	ControllerHandler
}

func New(sm session.SessionManager, notFound ControllerHandler) Router {
	return Router{
		sm,
		make([]UrlRoute, 0),
		notFound,
	}
}

func (router Router) findRoute (path string) UrlRoute {
	for _, v := range router.urlRoutes {
		pathRegExp := regexp.MustCompile(v.urlRegExp)

		if pathRegExp.MatchString(path) {
			return v
		}
	}
	return UrlRoute{
		handler: router.pageNotFoundController,
	}
}

func (router *Router) New(sm session.ISessionManager) {
	router.sessionManager = sm
}

func (router *Router) Route(w http.ResponseWriter, r *http.Request)  {
	urlPath := r.URL.Path
	route := router.findRoute(urlPath)
	params := make(map[string]string)
	pathItems := strings.Split(urlPath, "/")

	for k, v := range route.params {
		if len(pathItems) > v {
			params[k] = pathItems[v]
		}
	}

	urlOptions := &UrlOptions{
		 params,
	}

	logger.Log(logger.INFO, "Navigating to url = " + urlPath + " vs route = " +
        route.urlRegExp)

	routeHandler := route.handler
	routeHandler(w, r, *urlOptions, router.sessionManager)
}

func (router *Router) AddRoute(urlPattern string, pathHandler ControllerHandler) {

	params := make(map[string]int)
	pathRegExp := url.UrlPatternToRegExp(urlPattern)

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

	router.urlRoutes = append(router.urlRoutes, UrlRoute{
		urlRegExp: pathRegExp,
		handler: pathHandler,
		params: params,
	})
}
