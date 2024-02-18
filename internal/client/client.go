package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/AidanThomas/mercury/internal/log"
	"github.com/AidanThomas/mercury/internal/message"
)

type Message message.Message

var (
	user string
	id   string
)

func Start(addr string, usr string) {
	user = usr

	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	// Username handshake
	jMsg, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	var req Message
	if err := json.Unmarshal([]byte(jMsg), &req); err != nil {
		log.Errorf(err.Error())
		return
	}
	if req.Body != "USERNAME\n" {
		log.Errorf("Unexpected handshake, expected USERNAME, got: %s", req.Body)
		return
	}
	res, err := json.Marshal(Message{
		Body: user,
		User: user,
	})
	if err != nil {
		log.Errorf(err.Error())
	}
	Send(c, string(res))
	jMsg, err = bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	var confirm Message
	if err = json.Unmarshal([]byte(jMsg), &confirm); err != nil {
		log.Errorf(err.Error())
	}
	id = strings.TrimSpace(confirm.ConnId)

	out := make(chan Message)
	in := make(chan Message)

	go waitForIncoming(c, in)
	go waitForOutgoing(c, out)
	active := true

	for active {
		select {
		case msg := <-in:
			fmt.Print(">> " + msg.Body)
		case msg := <-out:
			if msg.Body == "STOP\n" {
				fmt.Println("TCP client exiting...")
				active = false
			}
			jMsg, err := json.Marshal(msg)
			if err != nil {
				log.Errorf(err.Error())
				return
			}
			Send(c, string(jMsg))
		}
	}
}

func waitForIncoming(c net.Conn, in chan Message) {
	for {
		data, _ := bufio.NewReader(c).ReadString('\n')
		var msg Message
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			log.Errorf(err.Error())
			return
		}
		in <- msg
	}
}

func waitForOutgoing(c net.Conn, out chan Message) {
	for {
		reader := bufio.NewReader(os.Stdin)
		body, err := reader.ReadString('\n')
		body = strings.TrimSpace(body)
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		out <- Message{
			Body:   body,
			ConnId: id,
			User:   user,
		}
	}
}

func Send(c net.Conn, msg string) error {
	_, err := c.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}
	return nil
}
