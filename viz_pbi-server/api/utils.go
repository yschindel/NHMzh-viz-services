package api

import (
	"strconv"
)

// StringToInt converts a string to an integer, returning 0 if the conversion fails
func StringToInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
