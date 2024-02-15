package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func waitForIncoming(c net.Conn, in chan string) {
	for {
		msg, _ := bufio.NewReader(c).ReadString('\n')
		in <- msg
		if msg == "STOP" {
			break
		}
	}
}

func waitForOutgoing(c net.Conn, out chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		out <- msg
		if msg == "STOP" {
			break
		}
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	out := make(chan string)
	in := make(chan string)

	go waitForIncoming(c, in)
	go waitForOutgoing(c, out)
	loop := true

	for loop {
		select {
		case msg := <-in:
			fmt.Print(">> " + msg)
		case msg := <-out:
			if msg == "STOP" {
				fmt.Println("TCP client exiting...")
				close(in)
				close(out)
				loop = false
				break
			}
			fmt.Fprintf(c, msg+"\n")
		}
	}
}
