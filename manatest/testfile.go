// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"errors"
	"fmt"
	"github.com/mattrmiller/go-mana-test/console"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// TestHeader is a structure to handle headers for a test.
type TestHeader struct {

	// Key, hold the key of the header.
	Key string `yaml:"key"`

	// Value, hold the value of the header.
	Value string `yaml:"value"`
}

// TestChecks is a structure to handle headers for a test.
type TestChecks struct {

	// Name, hold the name of this check for the test.
	Name string `yaml:"name"`

	// Check, hold the check for the test.
	Check string `yaml:"check"`

	// Value, hold the value of the test.
	Value string `yaml:"value"`
}

// TestFile is a structure to handle an individual test file.
type TestFile struct {

	// filePath, holds the path to the file.
	filePath string

	// Name, holds the name of the test.
	Name string `yaml:"name"`

	// Url, holds the url of the test.
	Url string `yaml:"url"`

	// Method, holds the http method of the test.
	Method string `yaml:"method"`

	// Index, hold the index of this test.
	Index int `yaml:"index"`

	// Headers, holds the header variables.
	Headers []TestHeader `yaml:"headers"`

	// Body, holds the test http body.
	Body interface{}

	// TestChecks, holds the checks variables.
	Checks []TestChecks `yaml:"checks"`
}

// Reads a test file.
func ReadTestFile(pathFile string) (*TestFile, error) {

	// Check if yml
	if filepath.Ext(strings.TrimSpace(pathFile)) != ".yml" {
		return nil, errors.New(fmt.Sprintf("Test file is not a `yml` file: %s", pathFile))
	}

	// Read file
	source, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to read test file at: %s", pathFile))
	}

	// Unmarshal yaml
	var testFile TestFile
	err = yaml.Unmarshal(source, &testFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid test file format at: %s", pathFile))
	}

	// Set path
	testFile.filePath = pathFile

	return &testFile, nil
}

// Validates a test file is in proper format.
func (testFile *TestFile) Validate() error {

	// Must have a name
	if len(testFile.Name) == 0 {
		return errors.New("Test file must have `name` field.")
	}

	// Must have a url
	if len(testFile.Url) == 0 {
		return errors.New("Test file must have `url` field.")
	}

	// Must have a method
	if len(testFile.Method) == 0 {
		return errors.New("Test file must have `method` field.")
	}
	if !ValidateMethod(&testFile.Method) {
		return errors.New("Test file has invalid `method` field.")
	}

	// Correct index
	if testFile.Index < 0 {
		testFile.Index = 0
	}

	// Convert body to json
	if testFile.Body != nil {
		body, err := ConvertYamlToJson(testFile.Body)
		if err != nil {
			return errors.New("Unable to unmarshal JSON body.")
		}
		testFile.Body = body
	}

	// Validate headers
	for _, header := range testFile.Headers {

		// -- Key
		if len(header.Key) == 0 {
			return errors.New("Test file header must have `key` field.")
		}

		// -- Value
		if len(header.Value) == 0 {
			return errors.New("Test file header must have `value` field.")
		}
	}

	// Validate checks
	for _, check := range testFile.Checks {

		// -- Name
		if len(check.Name) == 0 {
			return errors.New("Test file check must have `name` field.")
		}

		// -- Check
		if len(check.Check) == 0 {
			return errors.New("Test file check must have `check` field.")
		}
		check.Check = strings.ToLower(check.Check)
		if !ValidateCheck(&check.Check) {
			return errors.New("Test file check has an invalid `check` value.")
		}

		// -- Check
		if len(check.Value) == 0 {
			return errors.New("Test file check must have `value` field.")
		}
	}

	return nil
}

// Test runs the test
func (testFile *TestFile) Test(projFile *ProjectFile) error {

	// Replace headers global values
	var headers []TestHeader
	for _, header := range testFile.Headers {
		header.Value = ReplaceVars(header.Value, &projFile.Globals)
		headers = append(headers, header)
	}
	testFile.Headers = headers

	// Lets for the Url, with substitutions
	url := ReplaceVars(testFile.Url, &projFile.Globals)

	// Console
	console.Print(fmt.Sprintf("\tRunning test: %s...", testFile.Name))
	console.Print(fmt.Sprintf("\t\t%s: %s", testFile.Method, url))

	// Get client
	client := testFile.getRestyClient(projFile)

	// Set body
	if testFile.Body != nil && (testFile.Method == http.MethodPost || testFile.Method == http.MethodPut) {
		fmt.Print(ReplaceVars(testFile.Body.(string), &projFile.Globals))
		client.SetBody(ReplaceVars(testFile.Body.(string), &projFile.Globals))
	}

	// Run
	response, err := client.Execute(testFile.Method, url)
	if err != nil {
		return err
	}

	// Run tests
	console.Print("\t\tRunning checks...")
	err = RunChecks(&testFile.Checks, response)
	if err != nil {
		console.Print(fmt.Sprintf("\t\tFAIL: %s", err))
		console.PrintVerbose("")
		console.PrintVerbose(fmt.Sprintf("\t\t%s", string(response.Body())))
		console.PrintVerbose("")
	} else {
		console.Print("\t\tPASSED!")
	}

	return nil
}

// Gets the path of the test file.
func (testFile *TestFile) GetPath() string {
	return filepath.Dir(testFile.filePath)
}

// Gets the path to the test file.
func (testFile *TestFile) GetFilePath() string {
	return testFile.filePath
}

// Gets a resty client.
func (testFile *TestFile) getRestyClient(projFile *ProjectFile) *resty.Request {

	// Turn off debug
	resty.SetDebug(false)

	// Create client
	client := resty.NewRequest().
		SetContentLength(true)

	// Set headers
	for _, header := range testFile.Headers {
		client = client.SetHeader(header.Key, header.Value)
	}

	return client
}
