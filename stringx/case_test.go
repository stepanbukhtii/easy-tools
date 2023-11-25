package stringx

import (
	"reflect"
	"testing"
)

func TestToCamel(t *testing.T) {
	var tests = []struct {
		values []string
		expect string
	}{
		{values: []string{"first", "second"}, expect: "firstSecond"},
		{values: []string{"first", "second", "id", "http"}, expect: "firstSecondIDHTTP"},
		{values: []string{"first", "id", "http", "second"}, expect: "firstIDHTTPSecond"},
		{values: []string{"id", "first", "second"}, expect: "idFirstSecond"},
	}

	for _, test := range tests {
		result := ToCamel(test.values, true)
		if result != test.expect {
			t.Errorf("Get %#v, expected %#v", result, test.expect)
		}
	}
}

func TestToSnake(t *testing.T) {
	var tests = []struct {
		values []string
		expect string
	}{
		{values: []string{"first", "second"}, expect: "first_second"},
		{values: []string{"first", "second", "id", "http"}, expect: "first_second_id_http"},
	}

	for _, test := range tests {
		result := ToSnake(test.values)
		if result != test.expect {
			t.Errorf("Get %#v, expected %#v", result, test.expect)
		}
	}
}

func TestToSnakeAllUpper(t *testing.T) {
	var tests = []struct {
		values []string
		expect string
	}{
		{values: []string{"first", "second"}, expect: "FIRST_SECOND"},
		{values: []string{"first", "second", "id", "http"}, expect: "FIRST_SECOND_ID_HTTP"},
	}

	for _, test := range tests {
		result := ToSnakeAllUpper(test.values)
		if result != test.expect {
			t.Errorf("Get %#v, expected %#v", result, test.expect)
		}
	}
}

func TestToKebab(t *testing.T) {
	var tests = []struct {
		values []string
		expect string
	}{
		{values: []string{"first", "second"}, expect: "first-second"},
		{values: []string{"first", "second", "id", "http"}, expect: "first-second-id-http"},
	}

	for _, test := range tests {
		result := ToKebab(test.values)
		if result != test.expect {
			t.Errorf("Get %#v, expected %#v", result, test.expect)
		}
	}
}

func TestToKebabAllUpper(t *testing.T) {
	var tests = []struct {
		values []string
		expect string
	}{
		{values: []string{"first", "second"}, expect: "FIRST-SECOND"},
		{values: []string{"first", "second", "id", "http"}, expect: "FIRST-SECOND-ID-HTTP"},
	}

	for _, test := range tests {
		result := ToKebabAllUpper(test.values)
		if result != test.expect {
			t.Errorf("Get %#v, expected %#v", result, test.expect)
		}
	}
}

func TestSplitCamelCase(t *testing.T) {
	var tests = []struct {
		word   string
		expect []string
	}{
		{word: "firstSecond", expect: []string{"first", "second"}},
		{word: "firstSecondIDHTTP", expect: []string{"first", "second", "id", "http"}},
		{word: "firstIDHTTPSecond", expect: []string{"first", "id", "http", "second"}},
		{word: "idFirst", expect: []string{"id", "first"}},
		{word: "IDFirstSecond", expect: []string{"id", "first", "second"}},
	}

	for _, test := range tests {
		result := SplitCamelCase(test.word)
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf("Get %#v, expected %#v", result, test.expect)
		}
	}
}

func TestSplitSnakeCase(t *testing.T) {
	var tests = []struct {
		word   string
		expect []string
	}{
		{word: "first_second", expect: []string{"first", "second"}},
		{word: "first_second_id_http", expect: []string{"first", "second", "id", "http"}},
	}

	for _, test := range tests {
		result := SplitSnakeCase(test.word)
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf("Get %#v, expected %#v", result, test.expect)
		}
	}
}

func TestSplitKebabCase(t *testing.T) {
	var tests = []struct {
		word   string
		expect []string
	}{
		{word: "first-second", expect: []string{"first", "second"}},
		{word: "first-second-id-http", expect: []string{"first", "second", "id", "http"}},
	}

	for _, test := range tests {
		result := SplitKebabCase(test.word)
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf("Get %#v, expected %#v", result, test.expect)
		}
	}
}
