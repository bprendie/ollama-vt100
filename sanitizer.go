// ui/sanitizer.go

package ui

import (
	"strings"
	"unicode"
)

// ToASCII iterates over a string and replaces any non-ASCII characters
// with a placeholder '?'. This ensures compatibility with simple terminals.
func ToASCII(s string) string {
	var builder strings.Builder
	builder.Grow(len(s)) // Pre-allocate memory for efficiency

	for _, r := range s {
		if r > unicode.MaxASCII {
			builder.WriteRune('?')
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
