package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func stopWithFlag() {
	stop := false
	go func() {
		for !stop {
			pc, _, _, _ := runtime.Caller(0)
			funcName := runtime.FuncForPC(pc).Name()
			fmt.Printf("Work func %s\n", funcName)
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(2 * time.Second)
	stop = true
	time.Sleep(2 * time.Second)
}
func stopWithChannel() {
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Channel stopped")
				return
			default:
				fmt.Println("Channel running from function stopWithChannel")
				time.Sleep(1 * time.Second)
			}

		}
	}()
	time.Sleep(2 * time.Second)
	close(done)
	time.Sleep(2 * time.Second)
}
func stopWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context stopped")
				return
			default:
				fmt.Printf("Context running from function stopWithContext\n")
				time.Sleep(1 * time.Second)
			}

		}
	}(ctx)
	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(2 * time.Second)
}
func stopWithExit() {
	go func() {
		defer fmt.Println("Function stopWithExit stopped")

		fmt.Println("Function stopWithExit running")
		runtime.Goexit()
		time.Sleep(1 * time.Second)
	}()
	time.Sleep(1 * time.Second)
}

func stopWithTimeout() {
	go func() {
		timeout := time.After(2 * time.Second)
		for {
			select {
			case <-timeout:
				fmt.Println("Function stopWithTimeout stopped")
				return
			default:
				fmt.Println("Function stopWithTimeout running")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	time.Sleep(3 * time.Second)
}
func stopWithPanicRecover() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()

		fmt.Println("Function stopWithPanicRecover running")
		panic("Panic occurred")
	}()
	time.Sleep(1 * time.Second)
}

func main() {
	stopWithFlag()
	stopWithChannel()
	stopWithContext()
	stopWithExit()
	stopWithTimeout()
	stopWithPanicRecover()
}
