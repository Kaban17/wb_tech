package main

import (
	"fmt"
	"sync"
)

var cnt int
var mu sync.Mutex
var wg sync.WaitGroup

func inc() {
	mu.Lock()
	defer mu.Unlock()
	cnt++
	wg.Done()
}
func main() {
	fmt.Println("Hello, world from L1/18!")
	for range 100 {
		wg.Add(1)
		go inc()
	}
	wg.Wait()
	fmt.Println(cnt)
}
