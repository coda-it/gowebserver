package main

import (
    goWebServer "github.com/oskarszura/gowebserver"
)

var server goWebServer.WebServer

func main() {
    server.RunServer("9090")
}
