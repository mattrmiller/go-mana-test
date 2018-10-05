// Package console provides useful console functionality.
package console

// Imports
import (
	"github.com/fatih/color"
	"strings"
)

// Console colors
const (
	ColorWhite   = "white"
	ColorCyan    = "cyan"
	ColorRed     = "red"
	ColorGreen   = "green"
	ColorYellow  = "yellow"
	ColorMagenta = "magenta"
)

// Console structure handles all things console related.
type Console struct {

	// optColor determines if color console is used
	optColor bool
}

// NewConsole creates a new console structure with default options.
func NewConsole() *Console {

	return &Console{
		optColor: false,
	}
}

// formatMsg formats a console message.
func (cls *Console) formatMsg(msg string) string {
	return strings.Replace(msg, "\t", "  ", -1)
}

// SetOptColor sets color option.
func (cls *Console) SetOptColor(clr bool) {

	// Set
	cls.optColor = clr
}

// Print prints to console.
func (cls *Console) Print(msg string) {

	// Print
	cls.PrintColor(msg, ColorWhite)
}

// PrintError prints error to console.
func (cls *Console) PrintError(msg string) {

	// Print
	cls.PrintColor(msg, ColorRed)
}

// PrintColor prints to console, checking for color option.
func (cls *Console) PrintColor(msg string, clr string) {

	// Format message
	msg = cls.formatMsg(msg)

	// Print
	if cls.optColor {
		if clr == ColorCyan {
			color.Cyan(msg)
		} else if clr == ColorRed {
			color.Red(msg)
		} else if clr == ColorGreen {
			color.Green(msg)
		} else if clr == ColorYellow {
			color.Yellow(msg)
		} else if clr == ColorMagenta {
			color.Magenta(msg)
		} else {
			println(msg)
		}
	} else {
		println(msg)
	}
}
