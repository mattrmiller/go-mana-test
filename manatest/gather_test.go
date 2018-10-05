// Package manatest provides internal workings for go-mana-test.
package manatest

// Imports
import (
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test GatherTestFiles.
func TestGatherTestFiles(tst *testing.T) {

	// Gather files in valid folder
	testFiles, err := GatherTestFiles("../testproj/tests")
	brtesting.AssertEqual(tst, err, nil, "GatherTestFiles failed")
	brtesting.AssertEqual(tst, len(testFiles), 16, "GatherTestFiles failed to find the correct number of test files")

	// Gather files in valid folder
	_, err = GatherTestFiles("../testproj/tests/this-does-not-exist")
	brtesting.AssertNotEqual(tst, err, nil, "GatherTestFiles succeeded for a directory that does not exist")

}
