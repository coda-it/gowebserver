package gowebserver

import (
	"regexp"
	"net/http"
	"github.com/oskarszura/gowebserver/controllers"
)

var (
	routeFront = regexp.MustCompile(`^\/$`)
	routeNumbers = regexp.MustCompile(`\d`)
	routeApi = regexp.MustCompile(`^(\/api)(\/?\?{0}|\/?\?{1}.*)$`)
)

func Route(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	switch {
	case routeFront.MatchString(urlPath):
		controllers.ControllerFront(w, r)
	case routeNumbers.MatchString(urlPath):
		controllers.ControllerDigits(w, r)
	case routeApi.MatchString(urlPath):
		controllers.ControllerApi(w, r)
	default:
		controllers.Controller404(w, r)
	}
}