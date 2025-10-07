package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb_tech/l2_17/internal/telnet"
	types "wb_tech/l2_17/pkg"
)

var (
	host    string
	port    int
	timeout time.Duration
)

func main() {
	flag.StringVar(&host, "host", "localhost", "The host to connect to")
	flag.IntVar(&port, "port", 8080, "The port to connect to")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "The timeout for the connection")
	flag.Parse()
	config := types.NewConfig(host, port, timeout)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := telnet.Connect(ctx, config); err != nil {
		fmt.Fprintf(os.Stderr, "Client Error: %v\n", err)
		os.Exit(1)
	}

}
