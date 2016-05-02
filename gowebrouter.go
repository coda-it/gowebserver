package gowebserver

import (
	"regexp"
	"net/http"
	"io/ioutil"
	"html/template"
)

type Router struct {

}

var (
	rNum = regexp.MustCompile(`\d`)
	rAbc = regexp.MustCompile(`abc`)
)

type Page struct {
	Title string
	Body  template.HTML
}

func loadPage(title string) (*Page, error) {
	filename := "public/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: template.HTML(body)}, nil
}

func route(w http.ResponseWriter, r *http.Request) {
	switch {
	case rNum.MatchString(r.URL.Path):
		controllerDigits(w, r)
	case rAbc.MatchString(r.URL.Path):
		controllerAbc(w, r)
	default:
		controller404(w, r)
	}
}

func controller404(w http.ResponseWriter, r *http.Request) {
	p, _ := loadPage("404")
	t, _ := template.ParseFiles("public/view.html")
	t.Execute(w, p)
}

func controllerDigits(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Has digits"))
}

func controllerAbc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Has abc"))
}