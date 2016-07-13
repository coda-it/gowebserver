package gowebserver

import (
	"log"
	"os/exec"
)

func SidGenerator () string {
	sid, err := exec.Command("uuidgen").Output()

	if(err != nil) {
		log.Println(err)
		return "guest"
	}

	return string(sid)
}
