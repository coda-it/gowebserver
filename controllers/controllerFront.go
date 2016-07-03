package controllers

import (
	"net/http"
	"html/template"
	"github.com/oskarszura/gowebserver/models"
	"os"
	"path/filepath"
	"log"
)

func ControllerFront(w http.ResponseWriter, r *http.Request) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	p, _ := models.LoadPage("index")
	t, _ := template.ParseFiles(dir + "/public/view.html")
	t.Execute(w, p)
}
