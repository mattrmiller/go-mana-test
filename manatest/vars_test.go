// Package manatest provides internal workings for go-mana-test.
package manatest

// Imports
import (
	"fmt"
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
	str = ReplaceGlobalVars(str, &vars)
	brtesting.AssertEqual(tst, str, "Tucker is my dog. Fluffy is my cat.", "ReplaceGlobalVars failed for valid globals.")

}

// Test ReplaceRandomString.
func TestReplaceRandomString(tst *testing.T) {

	// Test 20 length
	str := ReplaceRandomString("{{rand.string.50}}")
	brtesting.AssertEqual(tst, len(str), 50, fmt.Sprintf("ReplaceRandomString failed to create a string of length 50: %s", str))

}

// Test ReplaceRandomStringLower.
func TestReplaceRandomStringLower(tst *testing.T) {

	// Test 20 length
	str := ReplaceRandomStringLower("{{rand.string.lower.40}}")
	brtesting.AssertEqual(tst, len(str), 40, fmt.Sprintf("ReplaceRandomStringLower failed to create a string of length 40: %s", str))

}

// Test ReplaceRandomStringUpper.
func TestReplaceRandomStringUpper(tst *testing.T) {

	// Test 20 length
	str := ReplaceRandomStringUpper("{{rand.string.upper.30}}")
	brtesting.AssertEqual(tst, len(str), 30, fmt.Sprintf("ReplaceRandomStringUpper failed to create a string of length 30: %s", str))

}

// Test ReplaceRandomNumber.
func TestReplaceRandomNumber(tst *testing.T) {

	// Test 20 length
	str := ReplaceRandomNumber("{{rand.num.1.9}}")
	brtesting.AssertEqual(tst, len(str), 1, fmt.Sprintf("ReplaceRandomNumber failed to create a number between 1 and 9: %s", str))

}
