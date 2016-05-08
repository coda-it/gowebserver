package controllers

import (
	"net/http"
	"html/template"
	"github.com/oskarszura/gowebserver/models"
)

func Controller404(w http.ResponseWriter, r *http.Request) {
	p, _ := models.LoadPage("404")
	t, _ := template.ParseFiles("public/view.html")
	t.Execute(w, p)
}
