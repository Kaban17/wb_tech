package main

import (
	"fmt"
	"time"
)

func main() {
	arr := []int{2, 4, 6, 8, 10}
	for _, num := range arr {
		go func(num int) {
			num = num * num
			fmt.Println(num)
		}(num)
	}
	time.Sleep(time.Millisecond * 100)
}
