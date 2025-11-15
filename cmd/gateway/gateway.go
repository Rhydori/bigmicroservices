package main

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rhydori/bigmicroservices/config"
	"github.com/rhydori/bigmicroservices/logs"
)

const (
	gateAddr  = "127.0.0.1:1000"
	loginAddr = "127.0.0.1:2000"
	chatAddr  = "127.0.0.1:3000"
)

type Gateway struct {
	id string

	gateAddr  string
	loginAddr string
	chatAddr  string

	loginConn net.Conn
	chatConn  net.Conn

	clients map[string]*Client
	mu      sync.RWMutex

	quitCh      chan struct{}
	broadcastCh chan Message
}

type Client struct {
	id     string
	conn   net.Conn
	sendCh chan Message
}

type Message struct {
	payload []byte
}

func main() {
	cfg := config.GetConfig()

	gateway := newGateway(cfg.GatewayAddr, cfg.LoginAddr, cfg.ChatAddr)
	gateway.startGateway()
}

func newGateway(gateAddr, loginAddr, chatAddr string) *Gateway {
	return &Gateway{
		id:        uuid.New().String(),
		gateAddr:  gateAddr,
		loginAddr: loginAddr,
		chatAddr:  chatAddr,
		clients:   make(map[string]*Client),

		quitCh:      make(chan struct{}),
		broadcastCh: make(chan Message, 32),
	}
}

func (g *Gateway) startGateway() {
	ln, err := net.Listen("tcp", g.gateAddr)
	if err != nil {
		logs.Fatal("Gateway start failed: %v", err)
	}
	defer ln.Close()

	logs.Info("Waiting for LoginServer...")
	g.loginConn = g.loginConnection()

	logs.Info("Waiting for ChatServer...")
	g.chatConn = g.chatConnection()

	logs.Info("Gateway listening at %s", ln.Addr())
	go g.acceptClientConn(ln)

	<-g.quitCh
	//close(g.msgCh)
}

func (g *Gateway) loginConnection() net.Conn {
	for {
		conn, err := net.Dial("tcp", g.loginAddr)
		if err == nil {
			logs.Info("Connected to LoginServer: %v", conn.RemoteAddr().String())
			return conn
		}
		time.Sleep(2 * time.Second)
	}
}

func (g *Gateway) chatConnection() net.Conn {
	for {
		conn, err := net.Dial("tcp", g.chatAddr)
		if err == nil {
			logs.Info("Connected to ChatServer: %v", conn.RemoteAddr().String())
			return conn
		}
		time.Sleep(2 * time.Second)
	}
}

func (g *Gateway) acceptClientConn(ln net.Listener) {
	logs.Info("Gateway accepting connections...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Warn("Gateway acceptConn Error: %v", err)
			continue
		}
		c := &Client{id: uuid.New().String(), conn: conn, sendCh: make(chan Message, 32)}

		g.mu.Lock()
		g.clients[c.id] = c
		g.mu.Unlock()

		logs.Debug("Connected: %s", c.conn.RemoteAddr())
		go g.handleRead(c)
		//go g.handleMessage(c)
		go g.broadcastMessage(c)
	}
}

func (g *Gateway) handleRead(c *Client) {
	defer func() {
		c.conn.Close()
		g.mu.Lock()
		delete(g.clients, c.id)
		g.mu.Unlock()
		close(c.sendCh)
		logs.Debug("Disconnected: %s", c.conn.RemoteAddr())
	}()

	buf := make([]byte, 2048)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			logs.Warn("Gateway handleRead Error: %s: %v", c.conn.RemoteAddr(), err)
			return
		}
		//		c.sendCh <- Message{
		//			payload: buf[:n],
		//		}
		g.broadcastCh <- Message{
			payload: buf[:n],
		}
	}
}

//func (g *Gateway) handleMessage(c *Client) {
//	for msg := range c.sendCh {
//		logs.Debug("handleMessage: %s: %s", c.conn.RemoteAddr(), msg.payload)
//	}
//}

func (g *Gateway) broadcastMessage(c *Client) {
	for msg := range g.broadcastCh {
		broadcastPayload := append([]byte(c.conn.RemoteAddr().String()+": "), msg.payload...)
		logs.Debug("%v", string(broadcastPayload))
		g.mu.Lock()
		for _, client := range g.clients {
			go func(cl *Client) {
				_, err := cl.conn.Write(broadcastPayload)
				if err != nil {
					logs.Warn("Broadcast Msg Error: %v", err)
				}

			}(client)
		}
		g.mu.Unlock()
	}
}
