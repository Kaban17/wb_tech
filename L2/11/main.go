package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
)

func sortStringByRunes(s string) string {
	runes := []rune(s)
	slices.Sort(runes)
	return string(runes)
}

func findAnagramSets(words []string) map[string][]string {
	tempMap := make(map[string][]string)

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		if !isLetterString(lowerWord) {
			continue
		}

		canonical := sortStringByRunes(lowerWord)

		tempMap[canonical] = append(tempMap[canonical], lowerWord)
	}

	result := make(map[string][]string)

	for _, anagrams := range tempMap {
		if len(anagrams) < 2 {
			continue
		}

		sort.Strings(anagrams)

		key := anagrams[0]
		result[key] = anagrams
	}

	return result
}

func isLetterString(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	anagrams := findAnagramSets(words)

	fmt.Println("Входные слова:", words)
	fmt.Println("Найденные множества анаграмм:")
	for key, value := range anagrams {
		fmt.Printf("- %q: %v\n", key, value)
	}
}
