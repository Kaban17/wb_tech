package main

import "fmt"

func intersection(a, b []int) []int {
	res := []int{}
	set := make(map[int]bool)
	for _, val := range a {
		set[val] = true
	}
	seen := make(map[int]bool)
	for _, val := range b {
		if set[val] && !seen[val] {
			res = append(res, val)
			seen[val] = true
		}
	}
	return res
}
func main() {
	fmt.Println("Hello, world from L1/11!")
	a := []int{1, 2, 3, 4, 5}
	b := []int{4, 5, 6, 7, 8}
	fmt.Println(intersection(a, b))
}
