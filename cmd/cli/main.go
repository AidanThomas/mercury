package main

import (
	"os"

	"github.com/AidanThomas/mercury/internal/client"
	"github.com/AidanThomas/mercury/internal/log"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		log.Errorf("Please provide host:port")
		return
	}

	CONNECT := arguments[1]
	client.Start(CONNECT)
}
