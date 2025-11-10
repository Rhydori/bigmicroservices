package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/rhydori/bigmicroservices/logs"
)

var (
	clients       = make(map[string]*Client)
	broadcastChan = make(chan []byte, 1024)
	mu            sync.RWMutex
)

type Client struct {
	ID     string
	Conn   net.Conn
	SendCh chan []byte
}

func Start() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		logs.Fatal("%v", err)
	}
	logs.Info("Listening TCP at %v", ln.Addr().String())

	go acceptConn(ln)
	go broadcaster()
}

func acceptConn(ln net.Listener) {
	logs.Info("Accepting connections...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Warn("%v", err)
			return
		}
		c := &Client{ID: uuid.New().String(), Conn: conn, SendCh: make(chan []byte, 32)}

		mu.Lock()
		clients[c.ID] = c
		mu.Unlock()

		logs.Info("Connected: %v", c.Conn.RemoteAddr().String())
		go handleRead(c)
		go handleWrite(c)
	}
}

func handleRead(c *Client) {
	defer func() {
		c.Conn.Close()
		close(c.SendCh)
		mu.Lock()
		delete(clients, c.ID)
		mu.Unlock()
		logs.Info("Disconnected: %v", c.Conn.RemoteAddr().String())
	}()
	buf := make([]byte, 1024)
	for {
		n, err := c.Conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			logs.Warn("%v: %v", c.Conn.RemoteAddr().String(), err)
			return
		}
		msg := []byte(fmt.Sprintf("%v: %v", c.Conn.RemoteAddr().String(), string(buf[:n])))
		logs.Debug("%v", string(msg))
		broadcastChan <- msg
	}
}

func handleWrite(c *Client) {
	for msg := range c.SendCh {
		_, err := c.Conn.Write(msg)
		if err != nil {
			logs.Warn("%v: %v", c.Conn.RemoteAddr().String(), err)
			return
		}
	}
}

func broadcaster() {
	for msg := range broadcastChan {
		mu.RLock()
		for _, client := range clients {
			go func(cl *Client) {
				_, err := client.Conn.Write(msg)
				if err != nil {
					logs.Warn("%v: %v", client.Conn.RemoteAddr().String(), err)
				}
			}(client)
		}
		mu.RUnlock()

	}
}
