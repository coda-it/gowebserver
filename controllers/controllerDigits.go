package controllers

import (
	"net/http"
)

func ControllerDigits(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Has digits"))
}
