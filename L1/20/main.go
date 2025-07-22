package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func reverseRunes(runes []rune, start, end int) {
	for start < end {
		runes[start], runes[end] = runes[end], runes[start]
		start++
		end--
	}
}

func reverseWordsInPlace(s string) string {
	runes := []rune(s)

	reverseRunes(runes, 0, len(runes)-1)

	start := 0
	for i := 0; i <= len(runes); i++ {
		if i == len(runes) || unicode.IsSpace(runes[i]) {
			reverseRunes(runes, start, i-1)
			start = i + 1
		}
	}

	return string(runes)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите строку: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	result := reverseWordsInPlace(input)
	fmt.Println("Перевёрнутая строка:", result)
}
