package types

import (
	"time"
)

type Config struct {
	Host    string
	Port    int
	Timeout time.Duration
}

func NewConfig(host string, port int, timeout time.Duration) *Config {
	return &Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}
}
