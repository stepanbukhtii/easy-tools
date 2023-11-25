package stringx

import (
	"testing"
)

func TestToPlural(t *testing.T) {
	var tests = []struct {
		word   string
		expect string
	}{
		{word: "car", expect: "cars"},
		{word: "day", expect: "days"},
		{word: "bus", expect: "buses"},
		{word: "city", expect: "cities"},
		{word: "hero", expect: "heroes"},
		{word: "userID", expect: "userIDs"},
	}

	for _, test := range tests {
		if test.expect != ToPlural(test.word) {
			t.Errorf("Get %#v, expected %#v", test.word, test.expect)
		}
	}
}
