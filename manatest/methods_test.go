// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test ValidateMethod.
func TestValidateMethod(tst *testing.T) {

	// Valid methods
	methods := []string{
		"GET",
		"get",
		"POST",
		"post",
		"PUT",
		"put",
		"DELETE",
		"delete",
		"OPTIONS",
		"options",
		"PATCH",
		"patch",
	}
	for _, method := range methods {
		valid := ValidateMethod(method)
		brtesting.AssertEqual(tst, valid, true, "ValidateMethod failed for valid methods")
	}

	// Invalid methods
	methods = []string{
		"sdfsfddg",
		"1323423423",
		"FLKDKDKS@",
		"ddddsdfsd",
	}
	for _, method := range methods {
		valid := ValidateMethod(method)
		brtesting.AssertEqual(tst, valid, false, "ValidateMethod failed for invalid methods")
	}
}
