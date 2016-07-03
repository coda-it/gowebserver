package controllers

import (
	"net/http"
	"html/template"
	"github.com/oskarszura/gowebserver/models"
	"os"
	"path/filepath"
	"log"
)

func Controller404(w http.ResponseWriter, r *http.Request) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	p, _ := models.LoadPage("404")
	t, _ := template.ParseFiles(dir + "/gpublic/view.html")
	t.Execute(w, p)
}
