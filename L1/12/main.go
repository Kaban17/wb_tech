package main

import "fmt"

func properSubset[T comparable](set []T) []T {
	res := make([]T, 0)
	seen := make(map[T]bool)
	for _, elem := range set {
		if !seen[elem] {
			res = append(res, elem)
			seen[elem] = true
		}
	}
	return res
}
func main() {
	set := []string{"cat", "dog", "cat", "dog", "cat", "dog", "bird", "bird"}
	fmt.Println(properSubset(set))
}
