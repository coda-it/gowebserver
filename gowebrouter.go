package gowebserver

import (
	"net/http"
)

func SetupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
}