// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"net/http"
	"strings"
)

// Methods
const (
	METHOD_GET     = http.MethodGet
	METHOD_POST    = http.MethodPost
	METHOD_PUT     = http.MethodPut
	METHOD_DELETE  = http.MethodDelete
	METHOD_OPTIONS = http.MethodOptions
	METHOD_PATCH   = http.MethodPatch
)

// Validates http method.
func ValidateMethod(method *string) bool {

	// Methods
	*method = strings.ToUpper(*method)
	switch *method {
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
