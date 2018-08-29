// Package console provides useful console functionality.
package console

import "strings"

// Global variables
var (

	// verboseMode, determines if console will work in verbose mode
	verboseMode = false
)

// SetVerboseMode Sets console in verbose mode
func SetVerboseMode() {

	// Set
	verboseMode = true
}

// Print Prints a console line
func Print(msg string) {
	msg = strings.Replace(msg, "\t", "  ", -1)
	println(msg)
}

// PrintVerbose Prints a console line that is verbose
func PrintVerbose(msg string) {
	if verboseMode {
		Print(msg)
	}
}
