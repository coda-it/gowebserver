package controllers

import (
	"net/http"
	"html/template"
	"github.com/oskarszura/gowebserver/models"
)

func ControllerFront(w http.ResponseWriter, r *http.Request) {
	p, _ := models.LoadPage("index")
	t, _ := template.ParseFiles("public/view.html")
	t.Execute(w, p)
}
