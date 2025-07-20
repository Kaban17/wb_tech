package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(id int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d received cancellation signal\n", id)
			return
		default:
			fmt.Printf("Worker %d is working\n", id)
			time.Sleep(time.Second)
		}
	}
}
func main() {
	fmt.Println("Programm is running. Press Ctrl+C to exit.")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	num := 3
	for i := range num {
		wg.Add(1)
		go worker(i, ctx, &wg)
	}
	go func() {
		sig := <-sigChan
		fmt.Printf("Received signal: %s\n", sig)
		cancel()
	}()
	wg.Wait()

}
