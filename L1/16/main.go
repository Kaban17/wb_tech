package main

import (
	"fmt"
)

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	pivot := arr[0]
	less := []int{}
	greater := []int{}

	for _, val := range arr[1:] {
		if val <= pivot {
			less = append(less, val)
		} else {
			greater = append(greater, val)
		}
	}

	return append(append(quickSort(less), pivot), quickSort(greater)...)
}

func main() {
	arr := []int{10, 5, 2, 3, 8, 9, 1, 4}
	sorted := quickSort(arr)
	fmt.Println("Отсортированный массив:", sorted)
}
