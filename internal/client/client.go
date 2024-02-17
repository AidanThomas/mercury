package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/AidanThomas/mercury/internal/log"
)

func Start(addr string) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	out := make(chan string)
	in := make(chan string)

	go waitForIncoming(c, in)
	go waitForOutgoing(c, out)
	active := true

	for active {
		select {
		case msg := <-in:
			fmt.Print(">> " + msg)
		case msg := <-out:
			if msg == "STOP\n" {
				fmt.Println("TCP client exiting...")
				active = false
			}
			fmt.Fprintf(c, msg+"\n")
		}
	}
}

func waitForIncoming(c net.Conn, in chan string) {
	for {
		msg, _ := bufio.NewReader(c).ReadString('\n')
		in <- msg
	}
}

func waitForOutgoing(c net.Conn, out chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		out <- msg
	}
}
