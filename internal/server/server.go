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

	in := make(chan string)
	out := make(chan string)

	go awaitConnection(l, in, out)

	for {
		select {
		case msg := <-in:
			fmt.Println(msg)
		default:
			continue
		}
	}
}

func awaitConnection(l net.Listener, in, out chan string) {
	c, err := l.Accept()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	go handleConnection(c, in, out)
	go awaitConnection(l, in, out)
}

func handleConnection(c net.Conn, in, out chan string) {
	go waitForMessage(c, out)
	loop := true
	for loop {
		select {
		case msg := <-out:
			if msg == "STOP" {
				loop = false
				break
			}
			in <- msg
			c.Write([]byte("SERVER\n"))
		default:
			continue
		}
	}
	c.Close()
}

func waitForMessage(c net.Conn, out chan string) {
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
		out <- msg
	}
	c.Close()
}
