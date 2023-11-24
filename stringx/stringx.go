// Package stringx provides helper functions for working with strings.
package stringx

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
)

var compactRE = regexp.MustCompile(`\s+`)

// Shorten shortens a string to a specified number of characters.
func Shorten(s string, maxlen int) string {
	var short = []rune(s)
	if len(short) > maxlen {
		short = short[:maxlen]
		short = append(short, []rune(" [truncated]")...)
	}
	return string(short)
}

// Compact replaces consecutive whitespaces with a single space.
func Compact(s string) string {
	return compactRE.ReplaceAllString(string(s), " ")
}

// RandString generates a random string.
// length must be even.
func RandString(length int) string {
	b := make([]byte, length/2)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
