package repository

import "strings"

func RepeatWithSep(s string, count int, sep string) string {
	return strings.TrimSuffix(strings.Repeat(s+sep, count), sep)
}
