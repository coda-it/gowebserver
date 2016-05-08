package gowebserver

import (
	"regexp"
	"net/http"
	"html/template"
)

var (
	routeNumbers = regexp.MustCompile(`\d`)
	routeApi = regexp.MustCompile(`api`)
)

func Route(w http.ResponseWriter, r *http.Request) {
	switch {
	case routeNumbers.MatchString(r.URL.Path):
		controllerDigits(w, r)
	case routeApi.MatchString(r.URL.Path):
		controllerApi(w, r)
	default:
		controller404(w, r)
	}
}

func controller404(w http.ResponseWriter, r *http.Request) {
	p, _ := LoadPage("404")
	t, _ := template.ParseFiles("public/view.html")
	t.Execute(w, p)
}

func controllerDigits(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Has digits"))
}