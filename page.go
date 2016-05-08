package gowebserver

import (
	"html/template"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  template.HTML
}

func LoadPage(title string) (*Page, error) {
	filename := "public/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: template.HTML(body)}, nil
}