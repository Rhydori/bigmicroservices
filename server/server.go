package server

import (
	"fmt"
	"io"
	"net"

	"github.com/bolsonarius/logs"
	"github.com/google/uuid"
)

var clients = make(map[string]*Client)

type Client struct {
	ID   string
	Conn net.Conn
}

func Start() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		logs.Fatal("%v", err)
	}
	logs.Info("Listening TCP at %v", ln.Addr().String())
	go acceptConn(ln)
}

func acceptConn(ln net.Listener) {
	logs.Info("Accepting connections...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Warn("%v", err)
			return
		}
		c := &Client{ID: uuid.New().String(), Conn: conn}
		clients[c.ID] = c
		logs.Info("Connected: %v", c.Conn.RemoteAddr().String())
		go readPackets(c)
	}
}

func readPackets(c *Client) {
	defer func() {
		c.Conn.Close()
		logs.Info("Disconnected: %v", c.Conn.RemoteAddr().String())
		delete(clients, c.ID)
	}()
	for {
		buf := make([]byte, 1024)
		n, err := c.Conn.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			logs.Warn("%v: %v", c.Conn.RemoteAddr().String(), err)
			return
		}
		logs.Info("%v: %v", c.Conn.RemoteAddr().String(), string(buf[:n]))
		broadcastPackets(c, buf[:n])
	}
}

func broadcastPackets(c *Client, msg []byte) {
	msg = append([]byte(fmt.Sprintf("%v: ", c.Conn.RemoteAddr().String())), msg...)
	for _, client := range clients {
		client.Conn.Write(msg)
	}
}
