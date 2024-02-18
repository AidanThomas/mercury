package server

import (
	"bufio"
	"net"
)

type Connection struct {
	Active bool
	Id     string
	Conn   net.Conn
	User   string
}

func (c *Connection) Send(msg string) error {
	_, err := c.Conn.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) GetMsg() (string, error) {
	return bufio.NewReader(c.Conn).ReadString('\n')
}

func (c *Connection) Close() {
	c.Conn.Close()
	c.Active = false
}
