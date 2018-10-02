// Package http provides HTTP functionality for tests.
package http

// Imports
import (
	"net/http"
	"strings"
)

// Methods
const ()

// ValidateMethod Validates http method.
func ValidateMethod(method *string) bool {

	// Methods
	*method = strings.ToUpper(*method)
	switch *method {
	case
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodTrace,
		http.MethodHead:
		return true
	}
	return false

}
