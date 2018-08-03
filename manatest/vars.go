// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"fmt"
	"strings"
)

// Replace global variables in a string
func ReplaceGlobalVars(str string, vars []ProjectGlobal) string {

	// Only replace if we have our context
	if strings.Contains(str, "{{") && len(vars) != 0 {
		for _, val := range vars {
			replace := fmt.Sprintf("{{globals.%s}}", val.Key)
			str = strings.Replace(str, replace, val.Value, -1)
		}
	}

	return str
}
