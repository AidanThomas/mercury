package main

import (
	"os"

	"github.com/AidanThomas/mercury/internal/log"
	"github.com/AidanThomas/mercury/internal/server"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		log.Errorf("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	server.Start(PORT)
}
