package main

import "fmt"

func swap(a, b int) (int, int) {
	a = a + b
	b = a - b
	a = a - b
	return a, b
}
func main() {
	fmt.Println(swap(3, 10))
}
