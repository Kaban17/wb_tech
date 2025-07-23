package main

import "fmt"

func main() {
	data := []int{1, 2, 3, 4, 5}
	fmt.Println(data)
	fmt.Println(data[4])
	data = remove(data, 2) // Remove the element at index 2
	fmt.Println(data)
	// fmt.Println(data[4])  panic: runtime error: index out of range [4] with length 4
}
func remove(slice []int, index int) []int {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}
