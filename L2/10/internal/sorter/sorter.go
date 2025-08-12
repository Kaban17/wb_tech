package sorter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"wb_tech/l2_10/internal/config"
)

var monthMap = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

// Sort sorts the given lines based on the provided configuration.
func Sort(lines []string, cfg *config.Config) ([]string, error) {
	if cfg.CheckSorted {
		if isSorted(lines, cfg) {
			return nil, nil
		}
		return nil, fmt.Errorf("input is not sorted")
	}

	sort.SliceStable(lines, func(i, j int) bool {
		return less(lines[i], lines[j], cfg)
	})

	if cfg.Unique {
		lines = unique(lines, cfg)
	}

	return lines, nil
}

func less(s1, s2 string, cfg *config.Config) bool {
	f1 := getField(s1, cfg)
	f2 := getField(s2, cfg)

	if cfg.IgnoreTrail {
		f1 = strings.TrimRight(f1, " ")
		f2 = strings.TrimRight(f2, " ")
	}

	var result bool
	switch {
	case cfg.Numeric:
		n1, _ := strconv.Atoi(f1)
		n2, _ := strconv.Atoi(f2)
		result = n1 < n2
	case cfg.MonthSort:
		m1 := monthMap[strings.ToLower(f1[:3])]
		m2 := monthMap[strings.ToLower(f2[:3])]
		result = m1 < m2
	case cfg.HumanNumeric:
		n1 := parseHumanNumeric(f1)
		n2 := parseHumanNumeric(f2)
		result = n1 < n2
	default:
		result = f1 < f2
	}

	if cfg.Reverse {
		return !result
	}
	return result
}

func getField(line string, cfg *config.Config) string {
	if cfg.Column <= 1 {
		return line
	}
	fields := strings.Split(line, "\t")
	if cfg.Column > len(fields) {
		return ""
	}
	return fields[cfg.Column-1]
}

func parseHumanNumeric(s string) int64 {
	s = strings.ToLower(s)
	multiplier := int64(1)
	if strings.HasSuffix(s, "k") {
		multiplier = 1024
		s = strings.TrimSuffix(s, "k")
	} else if strings.HasSuffix(s, "m") {
		multiplier = 1024 * 1024
		s = strings.TrimSuffix(s, "m")
	} else if strings.HasSuffix(s, "g") {
		multiplier = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "g")
	}
	val, _ := strconv.ParseInt(s, 10, 64)
	return val * multiplier
}

func isSorted(lines []string, cfg *config.Config) bool {
	for i := 1; i < len(lines); i++ {
		if less(lines[i], lines[i-1], cfg) {
			return false
		}
	}
	return true
}

func unique(lines []string, cfg *config.Config) []string {
	if len(lines) == 0 {
		return lines
	}
	uniqueLines := make([]string, 0, len(lines))
	uniqueLines = append(uniqueLines, lines[0])
	for i := 1; i < len(lines); i++ {
		s1 := lines[i-1]
		s2 := lines[i]
		// Two lines are not equal if s1 < s2 or s2 < s1.
		// If they are equal, we skip adding the line.
		if less(s1, s2, cfg) || less(s2, s1, cfg) {
			uniqueLines = append(uniqueLines, s2)
		}
	}
	return uniqueLines
}
