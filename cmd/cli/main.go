package main

import (
	"os"

	"github.com/AidanThomas/mercury/internal/client"
	"github.com/AidanThomas/mercury/internal/log"
	"github.com/AidanThomas/mercury/internal/ui"
)

func main() {
	arguments := os.Args
	if len(arguments) != 3 {
		log.Errorf("Use like mercury <host:port> <username>")
		return
	}

	CONNECT := arguments[1]
	USER := arguments[2]
	ui.Start()
	client.Start(CONNECT, USER)
}
