// ui/textwrapper.go

package ui

import (
	"fmt"
	"strings"
)

// TextWrapper is a stateful writer that wraps text to a specified width.
type TextWrapper struct {
	width       int
	currentLine string
}

// NewTextWrapper creates a new TextWrapper.
func NewTextWrapper(width int) *TextWrapper {
	return &TextWrapper{width: width}
}

// Write adds new text and prints any completed, wrapped lines.
func (w *TextWrapper) Write(text string) {
	// Add the new text to the buffer, handling any existing newlines in the chunk.
	w.currentLine += strings.ReplaceAll(text, "\n", " \n ")

	// Process the buffer until it's shorter than the target width.
	for strings.Contains(w.currentLine, "\n") || len(w.currentLine) > w.width {
		var lineToPrint string

		// Find the earliest newline character.
		newlineIndex := strings.Index(w.currentLine, "\n")

		if newlineIndex != -1 && newlineIndex <= w.width {
			// If a newline is within the line width, break there.
			lineToPrint = w.currentLine[:newlineIndex]
			w.currentLine = w.currentLine[newlineIndex+1:]
		} else if len(w.currentLine) > w.width {
			// If the line is too long, find the best place to wrap.
			breakPoint := strings.LastIndex(w.currentLine[:w.width], " ")
			if breakPoint != -1 {
				// Wrap at the last space.
				lineToPrint = w.currentLine[:breakPoint]
				w.currentLine = strings.TrimSpace(w.currentLine[breakPoint:])
			} else {
				// No space found, so we have to break the long word.
				lineToPrint = w.currentLine[:w.width]
				w.currentLine = w.currentLine[w.width:]
			}
		} else {
			// No wrap needed yet, break the loop.
			break
		}

		fmt.Println(lineToPrint)
	}
}

// Flush prints any remaining text in the buffer.
func (w *TextWrapper) Flush() {
	if len(w.currentLine) > 0 {
		fmt.Println(strings.TrimSpace(w.currentLine))
		w.currentLine = ""
	}
}
