package gowebserver

import (
	"fmt"
	"net/http"
	"log"
)

var router Router

type WebServer struct {

}

func (s *WebServer) RunServer(port string) {
	http.HandleFunc("/", route)

	fmt.Println("Setting up server on " + port + " port")
	fmt.Println("Listening...")

	err := http.ListenAndServe(":" + port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

