// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test ReplaceGlobalVars.
func TestReplaceGlobalVars(tst *testing.T) {

	// Variables
	vars := []ProjectGlobal{
		{
			Key:   "DOG_NAME",
			Value: "Tucker",
		},
		{
			Key:   "CAT_NAME",
			Value: "Fluffy",
		},
	}
	str := "{{globals.DOG_NAME}} is my dog. {{globals.CAT_NAME}} is my cat."
	str = ReplaceGlobalVars(str, vars)
	brtesting.AssertEqual(tst, str, "Tucker is my dog. Fluffy is my cat.", "ReplaceGlobalVars failed for valid globals.")

}
