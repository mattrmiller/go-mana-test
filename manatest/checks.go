// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"errors"
	"fmt"
	"github.com/mattrmiller/go-mana-test/console"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"strconv"
	"strings"
)

// Checks
const (
	CheckResCode     = "response.code"
	CheckResBodyJSON = "response.body.json"

	CheckReqBodyJSON = "request.body.json"
)

// ValidateCheck Validates check.
func ValidateCheck(check *string) bool {

	// Response code
	if strings.HasPrefix(*check, CheckResCode) {
		return true
	}

	// Response body json
	if strings.HasPrefix(*check, CheckResBodyJSON) {
		return true
	}

	// Request body json
	if strings.HasPrefix(*check, CheckReqBodyJSON) {
		return true
	}

	return false
}

// RunChecks Runs checks for a test and a project.
func RunChecks(checks *[]TestChecks, vars *[]ProjectGlobal, response *resty.Response) error {

	// Each check
	for _, check := range *checks {

		// -- Verbose
		console.PrintVerbose(fmt.Sprintf("\t\t\t- %s", check.Name))

		// -- Replace variables
		check.Value = ReplaceVars(check.Value, vars)

		// -- Check response code
		if strings.HasPrefix(check.Check, CheckResCode) {
			err := checkResponseCode(&check, response)
			if err != nil {
				return err
			}
		}

		// -- Check response body json
		if strings.HasPrefix(check.Check, CheckResBodyJSON) {
			err := checkResponseBodyJSON(&check, response)
			if err != nil {
				return err
			}
		}

		// -- Check request body json
		if strings.HasPrefix(check.Check, CheckReqBodyJSON) {
			err := checkRequestBodyJSON(&check, response)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// checkResponseCode Checks the response code.
func checkResponseCode(check *TestChecks, response *resty.Response) error {

	// Convert response code in value
	resCode, err := strconv.Atoi(check.Value)
	if err != nil {
		return errors.New("invalid numeric 'response.code' in check value")
	}

	// Check
	if resCode != response.StatusCode() {
		return fmt.Errorf("check '%s' wanted '%d' != received '%d'", check.Check, resCode, response.StatusCode())
	}

	return nil
}

// checkResponseBodyJSON Checks the response body json.
func checkResponseBodyJSON(check *TestChecks, response *resty.Response) error {

	// First make sure that response type was json
	if !strings.Contains(response.Header().Get(HeaderContentType), ContentTypeJSON) {
		return fmt.Errorf("response '%s' was not a proper content type '%s'", check.Check, ContentTypeJSON)
	}

	// Scrape the prefix off the selector
	jsonSel := strings.TrimPrefix(check.Check, fmt.Sprintf("%s.", CheckResBodyJSON))

	// Get the json
	jsonValue := gjson.Get(string(response.Body()), jsonSel)

	// Check
	if !jsonValue.Exists() {
		return fmt.Errorf("check '%s' wanted '%s' != received 'null'", check.Check, check.Value)
	}
	if jsonValue.Str != check.Value {
		return fmt.Errorf("check '%s' wanted '%s' != received '%s'", check.Check, check.Value, jsonValue.Str)
	}

	return nil
}

// checkRequestBodyJSON Checks the request body json.
func checkRequestBodyJSON(check *TestChecks, response *resty.Response) error {

	// Scrape the prefix off the selector
	jsonSel := strings.TrimPrefix(check.Check, fmt.Sprintf("%s.", CheckReqBodyJSON))

	// Get the json
	jsonValue := gjson.Get(response.Request.Body.(string), jsonSel)

	// Check
	if !jsonValue.Exists() {
		return fmt.Errorf("check '%s' wanted '%s' != received 'null'", check.Check, check.Value)
	}
	if jsonValue.Str != check.Value {
		return fmt.Errorf("check '%s' wanted '%s' != received '%s'", check.Check, check.Value, jsonValue.Str)
	}

	return nil
}
