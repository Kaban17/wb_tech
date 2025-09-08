package config

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func LoadConfig() *Config {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is starting on port %s", port)
	return &Config{Port: port}
}
