package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/AidanThomas/mercury/internal/log"
)

var connections []*Connection

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
	conn := Connection{
		Id:   uint(len(connections) + 1),
		Conn: c,
	}
	connections = append(connections, &conn)
	go waitForIncoming(conn, in)
	go waitForOutgoing(conn, out)
}

func waitForIncoming(c Connection, in chan string) {
	for {
		msg, err := c.GetMsg()
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		msg = strings.TrimSpace(msg)
		if msg == "STOP" {
			break
		}
		in <- msg
	}
	c.Close()
}

func waitForOutgoing(c Connection, out chan string) {
	for {
		msg := <-out
		if err := c.Send(msg); err != nil {
			log.Errorf(err.Error())
		}
	}
}
