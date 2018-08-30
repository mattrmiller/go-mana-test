// Package brstrings provides useful testing functions.
package brtesting

// Imports
import (
	"fmt"
	"testing"
)

// Assert equal.
func AssertEqual(tst *testing.T, a interface{}, b interface{}, message string) {

	// Equal?
	if a == b {
		return
	}

	// Fatal
	tst.Fatal(fmt.Sprintf("%s: %v != %v", message, a, b))
}

// Assert not equal.
func AssertNotEqual(tst *testing.T, a interface{}, b interface{}, message string) {

	// Equal?
	if a != b {
		return
	}

	// Fatal
	tst.Fatal(fmt.Sprintf("%s: %v == %v", message, a, b))
}
