package gowebserver

import (
	"fmt"
	"net/http"
	"log"
)

type WebServer struct {

}

func (s *WebServer) RunServer(port string) {
	http.HandleFunc("/", Route)

	fmt.Println("Setting up server on " + port + " port")
	fmt.Println("Listening...")

	err := http.ListenAndServe(":" + port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

