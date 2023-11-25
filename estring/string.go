package estring

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(s string) string {
	return cases.Title(language.English).String(s)
}

func FirstLower(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func FirstUpper(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Before returns substring before substr.
func Before(s string, substr string) string {
	pos := strings.Index(s, substr)
	if pos == -1 {
		return ""
	}
	return s[:pos]
}

// After returns substring after substr.
func After(s string, substr string) string {
	pos := strings.LastIndex(s, substr)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(substr)
	if adjustedPos >= len(s) {
		return ""
	}
	return s[adjustedPos:]
}

// Between returns substring between two strings.
func Between(s, start, end string) string {
	posFirst := strings.Index(s, start)
	if posFirst == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(start)

	posLast := strings.Index(s[posFirstAdjusted:], end)
	if posLast == -1 {
		return s[posFirstAdjusted:]
	}
	return s[posFirstAdjusted : posFirstAdjusted+posLast]
}

// BetweenIndex returns start and end indexes between two strings.
func BetweenIndex(s, start, end string) (int, int) {
	posFirst := strings.Index(s, start)
	if posFirst == -1 {
		return -1, -1
	}
	posFirstAdjusted := posFirst + len(start)

	posLast := strings.Index(s[posFirstAdjusted:], end)
	if posLast == -1 || posFirstAdjusted >= posLast {
		return -1, -1
	}
	return posFirstAdjusted, posLast
}

// SplitByIndex returns a string by index from the split string
func SplitByIndex(s, sep string, index int) string {
	parts := strings.Split(s, sep)
	if index >= len(parts) {
		return ""
	}
	return parts[index]
}

// SplitAfter returns a string after substr from the split string
func SplitAfter(s, sep, substr string) string {
	parts := strings.Split(s, sep)
	for i := range parts {
		if parts[i] == substr && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// AfterWord returns a word after substr from the split string
func AfterWord(s, substr string) string {
	return SplitAfter(s, " ", substr)
}

// HasPrefixes returns string contains any one of the prefixes.
func HasPrefixes(s string, prefix ...string) bool {
	for i := range prefix {
		if strings.HasPrefix(s, prefix[i]) {
			return true
		}
	}
	return false
}

// HasSuffixes returns string contains any one of the suffixes.
func HasSuffixes(s string, suffix ...string) bool {
	for i := range suffix {
		if strings.HasSuffix(s, suffix[i]) {
			return true
		}
	}
	return false
}

// ContainsAny returns true if the s contains any ony of the substr
func ContainsAny(s string, substr ...string) bool {
	for i := range substr {
		if strings.Contains(s, substr[i]) {
			return true
		}
	}
	return false
}

// HasStrings returns true if the s equal any one of the substr
func HasStrings(s string, substr ...string) bool {
	for i := range substr {
		if s == substr[i] {
			return true
		}
	}
	return false
}

// EqualFoldAny returns true if the s equal with case-insensitivity any one of the substr
func EqualFoldAny(s string, substr ...string) bool {
	for i := range substr {
		if strings.EqualFold(s, substr[i]) {
			return true
		}
	}
	return false
}
