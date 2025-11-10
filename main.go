package main

import (
	"github.com/bolsonarius/server"
)

func main() {
	server.Start()

	select {}
}
