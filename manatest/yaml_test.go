// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"github.com/mattrmiller/go-bedrock/brtesting"
	"testing"
)

// Test ConvertYamlToJson.
func TestConvertYamlToJson(tst *testing.T) {

	// Variables
	yaml := `
  name: John Doe
  email: johndoe@gmail.com
`

	_, err := ConvertYamlToJson(yaml)
	brtesting.AssertEqual(tst, err, nil, "TestConvertYamlToJson failed for error not nil.")
}
