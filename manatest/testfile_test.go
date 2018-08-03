// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test ReadTestFile.
func TestReadTestFile(tst *testing.T) {

	// Valid files
	files := []string{
		"../testproj/tests/valid/simple.yml",
	}
	for _, file := range files {
		testFile, err := ReadTestFile(file)
		brtesting.AssertEqual(tst, err, nil, fmt.Sprintf("ReadTestFile failed: %s", file))
		brtesting.AssertEqual(tst, testFile.Validate(), nil, fmt.Sprintf("ReadTestFile failed to validate a valid test file: %s", file))
	}

	// Invalid files
	files = []string{
		"../testproj/tests/invalid/noname.yml",
		"../testproj/tests/invalid/nourl.yml",
		"../testproj/tests/invalid/nomethod.yml",
	}
	for _, file := range files {
		testFile, err := ReadTestFile(file)
		brtesting.AssertEqual(tst, err, nil, fmt.Sprintf("ReadTestFile failed: %s", file))
		brtesting.AssertNotEqual(tst, testFile.Validate(), nil, fmt.Sprintf("ReadTestFile failed to validate a invalid test file: %s", file))
	}

	// Non yml file
	_, err := ReadTestFile("../testproj/tests/this-does-not-exist.txt")
	brtesting.AssertNotEqual(tst, err, nil, "ReadTestFile did not fail on a test file that dos not exist")

	// Invalid test file
	_, err = ReadTestFile("../testproj/tests/this-does-not-exist.yml")
	brtesting.AssertNotEqual(tst, err, nil, "ReadTestFile did not fail on a test file that dos not exist")

}
