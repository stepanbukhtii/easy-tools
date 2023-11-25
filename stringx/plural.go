package stringx

import (
	"strings"
)

var (
	sibilantWords = []string{"s", "x", "z", "ss", "ch", "sh"}
	vowelWords    = map[byte]struct{}{'a': {}, 'e': {}, 'i': {}, 'o': {}, 'u': {}}
)

func ToPlural(s string) string {
	word := strings.ToLower(s)

	if HasSuffixes(word, sibilantWords...) {
		return s + "es"
	}

	l := len(word)
	if l >= 2 {
		if _, ok := vowelWords[word[l-2]]; !ok {
			if word[l-1] == 'o' {
				return s + "es"
			}
			if word[l-1] == 'y' {
				return s[:l-1] + "ies"
			}
		}
	}

	return s + "s"
}
