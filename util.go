package qlog

import (
	"os"
	"strings"
)

// is .ini or dir
func filter(input string) bool {
	if isIni(input) {
		return true
	}
	if f, err := os.Stat(input); err == nil && f.IsDir() {
		return true
	}
	return false
}

// is .ini file
func isIni(in string) bool {
	return strings.HasSuffix(in, ".ini")
}
