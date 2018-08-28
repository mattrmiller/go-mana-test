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
	CHECK_RES_CODE      = "response.code"
	CHECK_RES_BODY_JSON = "response.body.json"

	CHECK_REQ_BODY_JSON = "request.body.json"
)

// Validates check.
func ValidateCheck(check *string) bool {

	// Response code
	if strings.HasPrefix(*check, CHECK_RES_CODE) {
		return true
	}

	// Response body json
	if strings.HasPrefix(*check, CHECK_RES_BODY_JSON) {
		return true
	}

	// Request body json
	if strings.HasPrefix(*check, CHECK_REQ_BODY_JSON) {
		return true
	}

	return false
}

// Runs checks for a test and a project.
func RunChecks(checks *[]TestChecks, vars *[]ProjectGlobal, response *resty.Response) error {

	// Each check
	for _, check := range *checks {

		// -- Verbose
		console.PrintVerbose(fmt.Sprintf("\t\t\t- %s", check.Name))

		// -- Replace variables
		check.Value = ReplaceVars(check.Value, vars)

		// -- Check response code
		if strings.HasPrefix(check.Check, CHECK_RES_CODE) {
			err := checkResponseCode(&check, response)
			if err != nil {
				return err
			}
		}

		// -- Check response body json
		if strings.HasPrefix(check.Check, CHECK_RES_BODY_JSON) {
			err := checkResponseBodyJson(&check, response)
			if err != nil {
				return err
			}
		}

		// -- Check request body json
		if strings.HasPrefix(check.Check, CHECK_REQ_BODY_JSON) {
			err := checkRequestBodyJson(&check, response)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Checks the response code.
func checkResponseCode(check *TestChecks, response *resty.Response) error {

	// Convert response code in value
	resCode, err := strconv.Atoi(check.Value)
	if err != nil {
		return errors.New("Invalid numeric 'response.code' in check value.")
	}

	// Check
	if resCode != response.StatusCode() {
		return errors.New(fmt.Sprintf("Check '%s' wanted '%d' != received '%d'", check.Check, resCode, response.StatusCode()))
	}

	return nil
}

// Checks the response body json.
func checkResponseBodyJson(check *TestChecks, response *resty.Response) error {

	// First make sure that response type was json
	if !strings.Contains(response.Header().Get(HEADER_CONTENT_TYPE), CONTENT_TYPE_JSON) {
		return errors.New(fmt.Sprintf("Response '%s' was not a proper content type '%s'", check.Check, CONTENT_TYPE_JSON))
	}

	// Scrape the prefix off the selector
	jsonSel := strings.TrimPrefix(check.Check, fmt.Sprintf("%s.", CHECK_RES_BODY_JSON))

	// Get the json
	jsonValue := gjson.Get(string(response.Body()), jsonSel)

	// Check
	if !jsonValue.Exists() {
		return errors.New(fmt.Sprintf("Check '%s' wanted '%s' != received 'null'", check.Check, check.Value))
	}
	if jsonValue.Str != check.Value {
		return errors.New(fmt.Sprintf("Check '%s' wanted '%s' != received '%s'", check.Check, check.Value, jsonValue.Str))
	}

	return nil
}

// Checks the request body json.
func checkRequestBodyJson(check *TestChecks, response *resty.Response) error {

	// Scrape the prefix off the selector
	jsonSel := strings.TrimPrefix(check.Check, fmt.Sprintf("%s.", CHECK_REQ_BODY_JSON))

	// Get the json
	jsonValue := gjson.Get(response.Request.Body.(string), jsonSel)

	// Check
	if !jsonValue.Exists() {
		return errors.New(fmt.Sprintf("Check '%s' wanted '%s' != received 'null'", check.Check, check.Value))
	}
	if jsonValue.Str != check.Value {
		return errors.New(fmt.Sprintf("Check '%s' wanted '%s' != received '%s'", check.Check, check.Value, jsonValue.Str))
	}

	return nil
}
