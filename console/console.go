// Package console provides useful console functionality.
package console

// Global variables
var (

	// verboseMode, determines if console will work in verbose mode
	verboseMode = false
)

// Sets console in verbose mode
func SetVerboseMode() {

	// Set
	verboseMode = true
}

// Prints a console line
func Print(msg string) {
	println(msg)
}

// Prints a console line that is verbose
func PrintVerbose(msg string) {
	if verboseMode {
		println(msg)
	}
}
