package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/AidanThomas/mercury/internal/log"
)

func Start(port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	in := make(chan string)  // Incoming messages
	out := make(chan string) // Outgoing messages

	go awaitConnection(l, in, out)

	for {
		select {
		case msg := <-in:
			fmt.Println(msg)

			// Respond
			out <- "SERVER\n"
		default:
			continue
		}
	}
}

func awaitConnection(l net.Listener, in, out chan string) {
	for {
		c, err := l.Accept()
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		go handleConnection(c, in, out)
	}
}

func handleConnection(c net.Conn, in, out chan string) {
	// Create a listener to wait for an incoming message
	go waitForIncoming(c, in)
	go waitForOutgoing(c, out)
}

func waitForIncoming(c net.Conn, in chan string) {
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		msg := strings.TrimSpace(string(netData))
		if msg == "STOP" {
			break
		}
		in <- msg
	}
	c.Close()
}

func waitForOutgoing(c net.Conn, out chan string) {
	for {
		msg := <-out
		c.Write([]byte(msg))
	}
}
