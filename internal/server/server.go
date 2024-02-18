package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/AidanThomas/mercury/internal/log"
)

type message struct {
	body   string
	connId uint
}

var connections []*Connection

func Start(port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	in := make(chan message)  // Incoming messages
	out := make(chan message) // Outgoing messages

	go awaitConnection(l, in, out)

	for {
		select {
		case msg := <-in:
			log.Infof("[Id: %d]: %s", msg.connId, msg.body)
			// Respond
			out <- message{
				body: "SERVER\n",
			}
		default:
			continue
		}
	}
}

func awaitConnection(l net.Listener, in, out chan message) {
	for {
		c, err := l.Accept()
		if err != nil {
			log.Errorf(err.Error())
		}

		go handleConnection(c, in, out)
	}
}

func handleConnection(c net.Conn, in, out chan message) {
	conn := Connection{
		Id:   uint(len(connections) + 1),
		Conn: c,
	}
	connections = append(connections, &conn)
	go waitForIncoming(conn, in)
	go waitForOutgoing(conn, out)
}

func waitForIncoming(c Connection, in chan message) {
	for {
		msg, err := c.GetMsg()
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		msg = strings.TrimSpace(msg)
		in <- message{
			body:   msg,
			connId: c.Id,
		}
	}
}

func waitForOutgoing(c Connection, out chan message) {
	for {
		msg := <-out
		if err := c.Send(msg.body); err != nil {
			log.Errorf(err.Error())
			return
		}
	}
}
