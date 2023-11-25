package stringx

import (
	"strings"
	"unicode"
)

type StyleCase int

const (
	NoneCase StyleCase = iota
	// CamelCase is CamelCase
	CamelCase
	// CamelCaseFirstLower is camelCase
	CamelCaseFirstLower
	// SnakeCase is snake_case
	SnakeCase
	// SnakeCaseAllUpper is SNAKE_CASE
	SnakeCaseAllUpper
	// KebabCase is kebab-case
	KebabCase
	// KebabCaseAllUpper is KEBAB-CASE
	KebabCaseAllUpper
)

var acronym = map[string]string{
	"id":   "ID",
	"ip":   "IP",
	"url":  "URL",
	"http": "HTTP",
}

func ConvertCase(s string, inputCase, outputCase StyleCase) string {
	if inputCase == outputCase {
		return s
	}

	var parts []string
	switch inputCase {
	case CamelCase:
		parts = SplitCamelCase(s)
	case CamelCaseFirstLower:
		if outputCase == CamelCase {
			return FirstUpper(s)
		}
		parts = SplitCamelCase(s)
	case SnakeCase, SnakeCaseAllUpper:
		parts = SplitSnakeCase(s)
	case KebabCase, KebabCaseAllUpper:
		parts = SplitKebabCase(s)
	default:
		return s
	}

	switch outputCase {
	case CamelCase:
		return ToCamel(parts, false)
	case CamelCaseFirstLower:
		return ToCamel(parts, true)
	case SnakeCase:
		return ToSnake(parts)
	case SnakeCaseAllUpper:
		return ToSnakeAllUpper(parts)
	case KebabCase:
		return ToKebab(parts)
	case KebabCaseAllUpper:
		return ToKebabAllUpper(parts)
	default:
		return s
	}
}

func ConvertStringToCase(s string, outputCase StyleCase) string {
	var parts []string
	if strings.Contains(s, "_") {
		parts = SplitSnakeCase(s)
	} else if strings.Contains(s, "-") {
		parts = SplitKebabCase(s)
	} else if strings.Contains(s, " ") || strings.Contains(s, "\t") {
		parts = strings.Fields(s)
	} else {
		parts = SplitCamelCase(s)
	}

	switch outputCase {
	case CamelCase:
		return ToCamel(parts, false)
	case CamelCaseFirstLower:
		return ToCamel(parts, true)
	case SnakeCase:
		return ToSnake(parts)
	case SnakeCaseAllUpper:
		return ToSnakeAllUpper(parts)
	case KebabCase:
		return ToKebab(parts)
	case KebabCaseAllUpper:
		return ToKebabAllUpper(parts)
	default:
		return s
	}
}

func ToCamel(values []string, firstWordLower bool) string {
	var b strings.Builder
	for i, v := range values {
		if firstWordLower && i == 0 {
			b.WriteString(v)
			continue
		}

		if word, ok := acronym[v]; ok {
			b.WriteString(word)
			continue
		}

		b.WriteString(Title(v))
	}
	return b.String()
}

func ToSnake(values []string) string {
	return strings.Join(values, "_")
}

func ToSnakeAllUpper(values []string) string {
	upperValues := make([]string, len(values))
	for i, value := range values {
		upperValues[i] = strings.ToUpper(value)
	}
	return strings.Join(upperValues, "_")
}

func ToKebab(values []string) string {
	return strings.Join(values, "-")
}

func ToKebabAllUpper(values []string) string {
	upperValues := make([]string, len(values))
	for i, value := range values {
		upperValues[i] = strings.ToUpper(value)
	}
	return strings.Join(upperValues, "-")
}

func SplitCamelCase(s string) []string {
	var parts []string
	var start int
	for i, r := range s {
		if !unicode.IsUpper(r) || start > i {
			continue
		}

		if start < i {
			parts = append(parts, strings.ToLower(s[start:i]))
			start = i
		}

		for key, value := range acronym {
			if len(s) < i+len(value) {
				continue
			}

			if s[i:i+len(value)] == value {
				parts = append(parts, key)
				start = i + len(value)
				break
			}
		}
	}
	if start != len(s) {
		parts = append(parts, strings.ToLower(s[start:]))
	}
	return parts
}

func SplitSnakeCase(s string) []string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
	}
	return parts
}

func SplitKebabCase(s string) []string {
	parts := strings.Split(s, "-")
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
	}
	return parts
}
