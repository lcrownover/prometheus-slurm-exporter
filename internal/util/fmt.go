package util

import "strings"

func RemoveWhitespace(s string) string {
	return strings.Join(strings.Fields(s), "")
}
