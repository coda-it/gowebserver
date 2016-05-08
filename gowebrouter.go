package gowebserver

import (
	"regexp"
	"net/http"
	"github.com/oskarszura/gowebserver/controllers"
)

var (
	routeNumbers = regexp.MustCompile(`\d`)
	routeApi = regexp.MustCompile(`^(\/api)(\/?\?{0}|\/?\?{1}.*)$`)
)

func Route(w http.ResponseWriter, r *http.Request) {
	switch {
	case routeNumbers.MatchString(r.URL.Path):
		controllers.ControllerDigits(w, r)
	case routeApi.MatchString(r.URL.Path):
		controllers.ControllerApi(w, r)
	default:
		controllers.Controller404(w, r)
	}
}