package gowebserver

import (
	"fmt"
	"net/http"
	"log"
)

func SetupServer(port string) {
	SetupRoutes()

	fmt.Println("Setting up server on " + port + " port")
	fmt.Println("Listening...")

	err := http.ListenAndServe(":" + port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}