package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/AidanThomas/mercury/internal/encoding"
	"github.com/AidanThomas/mercury/internal/log"
)

type message struct {
	body   string
	connId string
	user   string
}

var (
	connections []*Connection
	encoder     encoding.Encoder
)

func Start(port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	// Create new encoder
	salt := os.Getenv("ENCODING_SALT")
	minLength, err := strconv.Atoi(os.Getenv("ENCODING_MIN_LENGTH"))
	encoder = *encoding.NewEncoder(salt, minLength)

	in := make(chan message)  // Incoming messages
	out := make(chan message) // Outgoing messages

	go awaitConnection(l, in, out)

	for {
		select {
		case msg := <-in:
			log.Infof("[%s]: %s", msg.user, msg.body)
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
	id, err := encoder.GenerateNewId()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	conn := Connection{
		Id:   id,
		Conn: c,
	}
	// USERNAME handshake
	conn.Send("USERNAME\n")
	user, err := conn.GetMsg()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	conn.User = strings.TrimSpace(user)
	log.Infof("Client { Id: %s, User: %s } connected from %s", conn.Id, conn.User, conn.Conn.RemoteAddr().String())

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
			user:   c.User,
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
