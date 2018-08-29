// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test ValidateCheck.
func TestValidateCheck(tst *testing.T) {

	// Valid checks
	checks := []string{
		"response.code",
		"response.body.json",
	}
	for _, check := range checks {
		valid := ValidateCheck(&check)
		brtesting.AssertEqual(tst, valid, true, "ValidateCheck failed for valid checks")
	}

	// Invalid checks
	checks = []string{
		"sdfsfddg",
		"1323423423",
		"FLKDKDKS@",
		"ddddsdfsd",
	}
	for _, check := range checks {
		valid := ValidateCheck(&check)
		brtesting.AssertEqual(tst, valid, false, "ValidateCheck failed for invalid checks")
	}
}
