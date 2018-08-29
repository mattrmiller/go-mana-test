// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test ConvertYamlToJSON.
func TestConvertYamlToJSON(tst *testing.T) {

	// Variables
	yaml := `
  name: John Doe
  email: johndoe@gmail.com
`

	_, err := ConvertYamlToJSON(yaml)
	brtesting.AssertEqual(tst, err, nil, "TestConvertYamlToJSON failed for error not nil.")
}
