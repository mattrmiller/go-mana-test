// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test TestReadProjectFile.
func TestReadProjectFile(tst *testing.T) {

	// Valid files
	files := []string{
		"../testproj/projects/valid/noglobals.yml",
		"../testproj/projects/valid/simple.yml",
	}
	for _, file := range files {
		projFile, err := ReadProjectFile(file)
		brtesting.AssertEqual(tst, err, nil, fmt.Sprintf("ReadProjectFile failed: %s", file))
		brtesting.AssertEqual(tst, projFile.Validate(), nil,
			fmt.Sprintf("ReadProjectFile failed to validate a valid project file: %s", file))
	}

	// Invalid files
	files = []string{
		"../testproj/projects/invalid/noname.yml",
		"../testproj/projects/invalid/notests.yml",
	}
	for _, file := range files {
		projFile, err := ReadProjectFile(file)
		brtesting.AssertEqual(tst, err, nil, fmt.Sprintf("ReadProjectFile failed: %s", file))
		brtesting.AssertNotEqual(tst, projFile.Validate(), nil,
			fmt.Sprintf("ReadProjectFile failed to validate a invalid project file: %s", file))
	}

	// Non yml file
	_, err := ReadProjectFile("../testproj/projects/this-does-not-exist.txt")
	brtesting.AssertNotEqual(tst, err, nil, "ReadProjectFile did not fail on a project file that dos not exist")

	// Invalid project file
	_, err = ReadProjectFile("../testproj/projects/this-does-not-exist.yml")
	brtesting.AssertNotEqual(tst, err, nil, "ReadProjectFile did not fail on a project file that dos not exist")

}
