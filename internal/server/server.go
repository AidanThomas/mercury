package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/AidanThomas/mercury/internal/encoding"
	"github.com/AidanThomas/mercury/internal/log"
	"github.com/AidanThomas/mercury/internal/message"
)

type Message message.Message

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

	in := make(chan Message)  // Incoming messages
	out := make(chan Message) // Outgoing messages

	go awaitConnection(l, in)
	go waitForOutgoing(out)

	for {
		select {
		case msg := <-in:
			log.Infof("[%s]: %s", msg.User, msg.Body)
			// Respond
			out <- msg
		default:
			continue
		}
	}
}

func awaitConnection(l net.Listener, in chan Message) {
	for {
		c, err := l.Accept()
		if err != nil {
			log.Errorf(err.Error())
		}

		go handleConnection(c, in)
	}
}

func handleConnection(c net.Conn, in chan Message) {
	id, err := encoder.GenerateNewId()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	conn := Connection{
		Active: true,
		Id:     id,
		Conn:   c,
	}

	// USERNAME handshake
	req, err := json.Marshal(Message{
		Body: "USERNAME\n",
	})
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	conn.Send(string(req))
	jMsg, err := conn.GetMsg()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	var res Message
	if err = json.Unmarshal([]byte(jMsg), &res); err != nil {
		log.Errorf(err.Error())
	}
	conn.User = strings.TrimSpace(res.Body)
	confirm, err := json.Marshal(Message{
		Body:   "CONFIRM",
		ConnId: conn.Id,
		User:   conn.User,
	})
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	conn.Send(string(confirm))
	log.Infof("Client { Id: %s, User: %s } connected from %s", conn.Id, conn.User, conn.Conn.RemoteAddr().String())

	connections = append(connections, &conn)
	go waitForIncoming(&conn, in)
}

func waitForIncoming(c *Connection, in chan Message) {
	for {
		jMsg, err := c.GetMsg()
		if err != nil {
			if err.Error() == "EOF" {
				disconnectClient(c)
				return
			}
			log.Errorf(err.Error())
			return
		}

		var msg Message
		if err := json.Unmarshal([]byte(jMsg), &msg); err != nil {
			log.Errorf(err.Error())
			return
		}

		in <- msg
	}
}

func waitForOutgoing(out chan Message) {
	for {
		msg := <-out
		jMsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		for _, c := range connections {
			if msg.ConnId == c.Id {
				continue
			}

			if err := c.Send(string(jMsg)); err != nil {
				log.Errorf(err.Error())
			}
			log.Debugf("Sent { %s } to { Id: %s, User: %s }", jMsg, c.Id, c.User)
		}
	}
}

func disconnectClient(c *Connection) {
	c.Close()
	var index = 1
	for i, conn := range connections {
		if c.Id == conn.Id {
			index = i
			break
		}
	}
	connections = append(connections[:index], connections[index+1:]...)
	log.Infof("Client { Id: %s, User: %s } disconnected", c.Id, c.User)
}
