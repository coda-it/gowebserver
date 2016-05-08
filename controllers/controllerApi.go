package controllers

import (
"net/http"
"encoding/json"
)

type Model struct {
Name      string    `json:"name"`
Completed bool      `json:"completed"`
}

type Models []Model

func ControllerApi(w http.ResponseWriter, r *http.Request) {

w.Header().Set("Content-Type", "application/json; charset=UTF-8")

output := Models{
Model{Name: "Param 1"},
Model{Name: "Param 2"},
}

json.NewEncoder(w).Encode(output)

}
