package main

import "fmt"

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}
func double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()
	return out
}
func print(in <-chan int) {
	for n := range in {
		fmt.Println(n)
	}
}
func main() {
	c1 := gen(1, 2, 3, 4, 5, 6)
	c2 := double(c1)
	print(c2)
}
