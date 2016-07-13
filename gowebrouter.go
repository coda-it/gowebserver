package gowebserver

import (
	"regexp"
	"net/http"
	"time"
	"github.com/oskarszura/gowebserver/controllers"
)

var (
	routeFront = regexp.MustCompile(`^\/$`)
	routeNumbers = regexp.MustCompile(`\d`)
	routeApi = regexp.MustCompile(`^(\/api)(\/?\?{0}|\/?\?{1}.*)$`)
)

func Route(w http.ResponseWriter, r *http.Request) {
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
