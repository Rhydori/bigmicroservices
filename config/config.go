package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rhydori/bigmicroservices/logs"
)

type Config struct {
	ipAddr      string
	GatewayAddr string
	LoginAddr   string
	ChatAddr    string
}

func GetConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		logs.Fatal("Error loading .env: %v", err)
	} else {
		logs.Info(".env loaded successfully.")
	}

	ipAddr := os.Getenv("IP_ADDR")
	if ipAddr == "" {
		logs.Fatal(".env IP_ADDR is nil")
	}
	GatewayAddr := os.Getenv("GATEWAY_ADDR")
	if GatewayAddr == "" {
		logs.Fatal(".env GATEWAY_ADDR is nil")
	}
	LoginAddr := os.Getenv("LOGIN_ADDR")
	if LoginAddr == "" {
		logs.Fatal(".env LOGIN_ADDR is nil")
	}
	ChatAddr := os.Getenv("CHAT_ADDR")
	if ChatAddr == "" {
		logs.Fatal(".env CHAT_ADDR is nil")
	}

	return &Config{
		ipAddr:      ipAddr,
		GatewayAddr: ipAddr + ":" + GatewayAddr,
		LoginAddr:   ipAddr + ":" + LoginAddr,
		ChatAddr:    ipAddr + ":" + ChatAddr,
	}
}
