package main

import (
	"fmt"
	"time"
)

func sleep(d time.Duration) {
	start := time.Now()
	for time.Since(start) < d {
		continue

	}
}
func main() {
	fmt.Println("Hello, world from L1/25!")
	fmt.Println(time.Now())
	time.Sleep(2 * time.Second)
	fmt.Println(time.Now())
}
