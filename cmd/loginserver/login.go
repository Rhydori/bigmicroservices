package main

import (
	"net"

	"github.com/rhydori/bigmicroservices/logs"
)

func main() {
	server := NewLoginServer("127.0.0.1:2000")
	server.Start()

	select {}
}

type Server struct {
	lnAddr string
}

type Message struct {
	From    string
	payLoad []byte
}

func NewLoginServer(lnAddr string) *Server {
	return &Server{
		lnAddr: lnAddr,
	}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", s.lnAddr)
	if err != nil {
		logs.Fatal("Login Start fatal: %v", err)
	}
	logs.Info("Login listening at %v", ln.Addr().String())
	go acceptConn(ln)
}

func acceptConn(ln net.Listener) {
	logs.Info("Login accepting connections...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Error("Accept Connection Error: %v", err)
			continue
		}
		logs.Info("Connected: %v", conn.RemoteAddr().String())
		handleRead(conn)
	}
}

func handleRead(conn net.Conn) {
	defer conn.Close()
	for {
		msg := &Message{From: conn.RemoteAddr().String(), payLoad: make([]byte, 512)}
		_, err := conn.Read(msg.payLoad)
		if err != nil {
			logs.Error("Read Connection Error: %v", err)
			continue
		}
		logs.Debug("%v: %v", msg.From, msg.payLoad)
	}
}
