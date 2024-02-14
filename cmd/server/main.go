package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var count = 0

func awaitConnection(l net.Listener, m chan string) {
	c, err := l.Accept()
	if err != nil {
		panic(err)
	}
	go handleConnection(c, m)
	count++
	go awaitConnection(l, m)
}

func handleConnection(c net.Conn, m chan string) {
	n := make(chan string)
	go waitForMessage(c, n)
	loop := true
	for loop {
		select {
		case msg := <-n:
			if msg == "STOP" {
				loop = false
				break
			}
			m <- msg
			c.Write([]byte("SERVER\n"))
		default:
			continue
		}
	}
	c.Close()
}

func waitForMessage(c net.Conn, n chan string) {
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			panic(err)
		}

		msg := strings.TrimSpace(string(netData))
		if msg == "STOP" {
			break
		}
		n <- msg
	}
	c.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	m := make(chan string)
	go awaitConnection(l, m)

	for {
		select {
		case msg := <-m:
			fmt.Println(msg)
		default:
			continue
		}
	}
}
