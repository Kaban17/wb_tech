package main

import (
	"fmt"
	"strings"
)

func HasUniqueChars(s string) bool {
	s = strings.ToLower(s)
	seen := make(map[rune]bool)
	for _, char := range s {
		if char == ' ' {
			continue
		}
		if seen[char] {
			return false
		}
		seen[char] = true
	}
	return true
}
func main() {
	fmt.Println(HasUniqueChars("He l f g п "))
}
