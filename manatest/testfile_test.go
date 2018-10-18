// Package manatest provides internal workings for go-mana-test.
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
		"../testproj/tests/valid/cache.yml",
		"../testproj/tests/valid/checks.yml",
		"../testproj/tests/valid/headers.yml",
		"../testproj/tests/valid/nochecks.yml",
		"../testproj/tests/valid/nochecks.yml",
		"../testproj/tests/valid/params.yml",
		"../testproj/tests/valid/req_body.yml",
		"../testproj/tests/valid/simple.yml",
		"../testproj/tests/valid/vars.yml",
	}
	for _, file := range files {
		testFile, err := ReadTestFile(file)
		brtesting.AssertEqual(tst, err, nil, fmt.Sprintf("ReadTestFile failed: %s", file))
		brtesting.AssertEqual(tst, testFile.Validate(), nil, fmt.Sprintf("ReadTestFile failed to validate a valid test file: %s", file))
	}

	// Invalid files
	files = []string{
		"../testproj/tests/invalid/invalidcache1.yml",
		"../testproj/tests/invalid/invalidcache2.yml",
		"../testproj/tests/invalid/invalidcheck1.yml",
		"../testproj/tests/invalid/invalidcheck2.yml",
		"../testproj/tests/invalid/invalidheaders.yml",
		"../testproj/tests/invalid/invalidparams.yml",
		"../testproj/tests/invalid/nomethod.yml",
		"../testproj/tests/invalid/noname.yml",
		"../testproj/tests/invalid/nourl.yml",
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
