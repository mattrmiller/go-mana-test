// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"strings"
)

// Methods
const (
	METHOD_GET     = "GET"
	METHOD_POST    = "POST"
	METHOD_PUT     = "PUT"
	METHOD_DELETE  = "DELETE"
	METHOD_OPTIONS = "OPTIONS"
	METHOD_PATCH   = "PATCH"
)

// Validates http method.
func ValidateMethod(method string) bool {

	// Methods
	method = strings.ToUpper(method)
	switch method {
	case
		METHOD_GET,
		METHOD_POST,
		METHOD_PUT,
		METHOD_DELETE,
		METHOD_OPTIONS,
		METHOD_PATCH:
		return true
	}
	return false

}
