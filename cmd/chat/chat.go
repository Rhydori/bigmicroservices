package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/rhydori/bigmicroservices/logs"
)

func main() {
	server := NewChatServer("127.0.0.1:3000")
	server.Start()

	select {}
}

type Server struct {
	lnAddr        string
	clients       map[string]*Client
	broadcastChan chan []byte
	mu            sync.RWMutex
}

type Client struct {
	ID     string
	Conn   net.Conn
	SendCh chan []byte
}

type Message struct {
	payLoad []byte
}

func NewChatServer(lnAddr string) *Server {
	return &Server{
		lnAddr:        lnAddr,
		clients:       make(map[string]*Client),
		broadcastChan: make(chan []byte),
	}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", s.lnAddr)
	if err != nil {
		logs.Fatal("%v", err)
	}
	logs.Info("Chat listening at %v", ln.Addr().String())

	go s.acceptConn(ln)
	go s.broadcaster()
}

func (s *Server) acceptConn(ln net.Listener) {
	logs.Info("Chat accepting connections...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Warn("%v", err)
			continue
		}
		c := &Client{ID: uuid.New().String(), Conn: conn, SendCh: make(chan []byte, 32)}

		s.mu.Lock()
		s.clients[c.ID] = c
		s.mu.Unlock()

		logs.Info("Connected: %v", c.Conn.RemoteAddr().String())
		go s.handleRead(c)
		go handleWrite(c)
	}
}

func (s *Server) handleRead(c *Client) {
	defer func() {
		c.Conn.Close()
		close(c.SendCh)
		s.mu.Lock()
		delete(s.clients, c.ID)
		s.mu.Unlock()
		logs.Info("Disconnected: %v", c.Conn.RemoteAddr().String())
	}()
	msg := &Message{payLoad: make([]byte, 1024)}
	for {
		n, err := c.Conn.Read(msg.payLoad)
		if err != nil {
			if err == io.EOF {
				return
			}
			logs.Warn("%v: %v", c.Conn.RemoteAddr().String(), err)
			return
		}
		msg := []byte(fmt.Sprintf("%v: %v", c.Conn.RemoteAddr().String(), string(msg.payLoad[:n])))
		logs.Debug("%v", string(msg))
		s.broadcastChan <- msg
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

func (s *Server) broadcaster() {
	for msg := range s.broadcastChan {
		s.mu.RLock()
		for _, client := range s.clients {
			go func(cl *Client) {
				_, err := client.Conn.Write(msg)
				if err != nil {
					logs.Warn("%v: %v", client.Conn.RemoteAddr().String(), err)
				}
			}(client)
		}
		s.mu.RUnlock()

	}
}
