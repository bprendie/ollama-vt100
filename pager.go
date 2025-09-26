// ui/pager.go

package ui

import (
	"bufio"
	"fmt"
	"strings"
)

// Pager is a stateful writer that wraps text and pauses after a set number of lines.
type Pager struct {
	width        int
	height       int
	linesPrinted int
	currentLine  string
	reader       *bufio.Reader // For reading user input to continue
}

// NewPager creates a new Pager.
func NewPager(width, height int, reader *bufio.Reader) *Pager {
	return &Pager{
		width:  width,
		height: height,
		reader: reader,
	}
}

// Write adds new text, printing wrapped lines and pausing if the page height is reached.
func (p *Pager) Write(text string) {
	p.currentLine += strings.ReplaceAll(text, "\n", " \n ")

	for strings.Contains(p.currentLine, "\n") || len(p.currentLine) > p.width {
		// Check if we need to pause before printing the next line.
		if p.linesPrinted >= p.height {
			p.pause()
		}

		var lineToPrint string
		newlineIndex := strings.Index(p.currentLine, "\n")

		if newlineIndex != -1 && newlineIndex <= p.width {
			lineToPrint = p.currentLine[:newlineIndex]
			p.currentLine = p.currentLine[newlineIndex+1:]
		} else if len(p.currentLine) > p.width {
			breakPoint := strings.LastIndex(p.currentLine[:p.width], " ")
			if breakPoint != -1 {
				lineToPrint = p.currentLine[:breakPoint]
				p.currentLine = strings.TrimSpace(p.currentLine[breakPoint:])
			} else {
				lineToPrint = p.currentLine[:p.width]
				p.currentLine = p.currentLine[p.width:]
			}
		} else {
			break
		}

		fmt.Println(lineToPrint)
		p.linesPrinted++
	}
}

// Flush prints any remaining text in the buffer.
func (p *Pager) Flush() {
	if len(p.currentLine) > 0 {
		if p.linesPrinted >= p.height {
			p.pause()
		}
		fmt.Println(strings.TrimSpace(p.currentLine))
		p.linesPrinted++
		p.currentLine = ""
	}
}

// pause prints a "More" prompt and waits for the user to press Enter.
func (p *Pager) pause() {
	fmt.Print("-- More --")
	p.reader.ReadString('\n') // Wait for Enter

	// Use carriage return to clear the "-- More --" line
	fmt.Print("\r          \r")

	p.linesPrinted = 0 // Reset line count for the next page
}
