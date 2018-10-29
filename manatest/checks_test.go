// Package manatest provides internal workings for go-mana-test.
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
		"response.body.json.{{cache.something}}",
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
		"response.{{cache.something}}.json",
	}
	for _, check := range checks {
		valid := ValidateCheck(&check)
		brtesting.AssertEqual(tst, valid, false, "ValidateCheck failed for invalid checks")
	}
}
