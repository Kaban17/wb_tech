package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func unPack(s string) (string, error) {
	// Обработка пустой строки
	if s == "" {
		return "", nil
	}

	result := make([]rune, 0)
	i := 0

	if unicode.IsDigit(rune(s[0])) {
		return "", errors.New("invalid string: starts with a digit")
	}

	for i < len(s) {
		r := rune(s[i])

		if r == '\\' {
			if i+1 >= len(s) {
				return "", errors.New("invalid escape sequence at the end of the string")
			}
			escapedRune := rune(s[i+1])
			result = append(result, escapedRune)
			i += 2
			continue
		}

		if unicode.IsDigit(r) {
			if len(result) == 0 {
				return "", errors.New("invalid string: digit without preceding character")
			}

			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", errors.New("failed to parse digit")
			}

			lastRune := result[len(result)-1]

			result = result[:len(result)-1]

			for j := 0; j < count; j++ {
				result = append(result, lastRune)
			}
			i++
			continue
		}

		result = append(result, r)
		i++
	}

	return string(result), nil
}

func main() {
	res, err := unPack("q4w")
	if err != nil {
		fmt.Printf(" err %s", err)
	}
	fmt.Printf("Unpacked string: %s", res)
}
