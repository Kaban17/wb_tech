package main

import (
	"fmt"
	"time"
)

func main() {
	const N = 5
	dataChan := make(chan int)
	done := make(chan struct{})
	go func() { // отправляем данные
		i := 0
		for {
			select {
			case <-done:
				fmt.Println("Sender done")
				return
			case dataChan <- i:
				i++
				time.Sleep(1 * time.Second)
			}
		}
	}()
	go func() { // получаем данные
		for val := range dataChan {
			fmt.Printf("Data recieved %d\n", val)
		}
	}()
	<-time.After(N * time.Second)
	close(done)
	close(dataChan)
	time.Sleep(500 * time.Millisecond)
}
