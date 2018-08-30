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

// TestChecks is a structure to handle checks for a test.
type TestChecks struct {

	// Name, hold the name of this check for the test.
	Name string `yaml:"name"`

	// Check, hold the check for the test.
	Check string `yaml:"check"`

	// Value, hold the value of the test.
	Value string `yaml:"value"`
}

// TestCache is a structure to handle cache for a test.
type TestCache struct {

	// Name, hold the name of this check for the test.
	Name string `yaml:"name"`

	// Value, hold the value of the test.
	Value string `yaml:"value"`
}

// TestFile is a structure to handle an individual test file.
type TestFile struct {

	// filePath, holds the path to the file.
	filePath string

	// Name, holds the name of the test.
	Name string `yaml:"name"`

	// Index, hold the index of this test.
	Index int `yaml:"index"`

	// URL, holds the url of the test.
	URL string `yaml:"url"`

	// Method, holds the http method of the test.
	RequestMethod string `yaml:"request.method"`

	// Headers, holds the header variables.
	RequestHeaders []TestHeader `yaml:"request.headers"`

	// Body, holds the test http body.
	ReqBody interface{} `yaml:"request.body"`

	// TestChecks, holds the checks variables.
	Checks []TestChecks `yaml:"checks"`

	// Cache, holds the cache variables.
	Cache []TestCache `yaml:"cache"`
}

// ReadTestFile Reads a test file.
func ReadTestFile(pathFile string) (*TestFile, error) {

	// Check if yml
	if filepath.Ext(strings.TrimSpace(pathFile)) != ".yml" {
		return nil, fmt.Errorf("test file is not a 'yml' file: %s", pathFile)
	}

	// Read file
	source, err := ioutil.ReadFile(pathFile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("unable to read test file at: %s", pathFile)
	}

	// Unmarshal yaml
	var testFile TestFile
	err = yaml.Unmarshal(source, &testFile)
	if err != nil {
		return nil, fmt.Errorf("invalid test file format at: %s", pathFile)
	}

	// Set path
	testFile.filePath = pathFile

	return &testFile, nil
}

// Validate Validates a test file is in proper format.
func (testFile *TestFile) Validate() error {

	// Must have a name
	if len(testFile.Name) == 0 {
		return errors.New("test file must have 'name' field")
	}

	// Must have a url
	if len(testFile.URL) == 0 {
		return errors.New("test file must have 'url' field")
	}

	// Must have a method
	if len(testFile.RequestMethod) == 0 {
		return errors.New("test file must have 'method' field")
	}
	if !ValidateMethod(&testFile.RequestMethod) {
		return errors.New("test file has invalid 'method' field")
	}

	// Correct index
	if testFile.Index < 0 {
		testFile.Index = 0
	}

	// Convert request body to json
	if testFile.ReqBody != nil {
		body, err := ConvertYamlToJSON(testFile.ReqBody)
		if err != nil {
			return errors.New("unable to unmarshal JSON body")
		}
		testFile.ReqBody = body
	}

	// Validate headers
	for _, header := range testFile.RequestHeaders {

		// -- Key
		if len(header.Key) == 0 {
			return errors.New("test file header must have 'key' field")
		}

		// -- Value
		if len(header.Value) == 0 {
			return errors.New("test file header must have 'value' fieldt")
		}
	}

	// Validate checks
	for _, check := range testFile.Checks {

		// -- Name
		if len(check.Name) == 0 {
			return errors.New("test file check must have 'name' fieldt")
		}

		// -- Check
		if len(check.Check) == 0 {
			return errors.New("test file check must have 'check' field")
		}
		check.Check = strings.ToLower(check.Check)
		if !ValidateCheck(&check.Check) {
			return fmt.Errorf("test file check has an invalid 'check' field: `%s'", check.Check)
		}

		// -- Check
		if len(check.Value) == 0 {
			return errors.New("test file check must have 'value' field")
		}
	}

	// Validate cache
	for _, cache := range testFile.Cache {

		// -- Name
		if len(cache.Name) == 0 {
			return errors.New("test file cache must have 'name' field")
		}

		// -- Value
		if len(cache.Value) == 0 {
			return errors.New("test file cache must have 'value' field")
		}
		cache.Value = strings.ToLower(cache.Value)
		if !ValidateCacheValue(&cache.Value) {
			return fmt.Errorf("test file cache has an invalid 'value' field: '%s'", cache.Value)
		}
	}

	return nil
}

// Test Runs the test.
func (testFile *TestFile) Test(projFile *ProjectFile) bool {

	// Replace headers global values
	var headers []TestHeader
	for _, header := range testFile.RequestHeaders {
		header.Value = ReplaceVars(header.Value, &projFile.Globals)
		headers = append(headers, header)
	}
	testFile.RequestHeaders = headers

	// Lets for the URL, with substitutions
	url := ReplaceVars(testFile.URL, &projFile.Globals)

	// Console
	console.Print(fmt.Sprintf("Running test: %s...", testFile.Name))
	console.Print(fmt.Sprintf("\t%s: %s", testFile.RequestMethod, url))

	// Get client
	client := testFile.getRestyClient()

	// Set body
	if testFile.ReqBody != nil && testFile.RequestMethod != http.MethodTrace {
		client.SetBody(ReplaceVars(testFile.ReqBody.(string), &projFile.Globals))
	}

	// Run
	console.Print("\tRunning request...")
	response, err := client.Execute(testFile.RequestMethod, url)
	if err != nil {
		console.Print(fmt.Sprintf("\tFAIL: %s", err))
		return false
	}

	// Save cache
	err = SaveCacheFromResponse(&testFile.Cache, response)
	if err != nil {
		console.PrintVerbose("")
		console.PrintVerbose(fmt.Sprintf("\tError saving cache: '%s'", err))
		console.PrintVerbose("")
		console.PrintVerbose(fmt.Sprintf("\t%s", string(response.Body())))
		console.PrintVerbose("")
		return false
	}

	// Run tests
	console.Print("\tRunning checks...")
	err = RunChecks(&testFile.Checks, &projFile.Globals, response)
	if err != nil {
		console.Print(fmt.Sprintf("\tFAIL: %s", err))
		console.PrintVerbose("")
		console.PrintVerbose(fmt.Sprintf("\t%s", string(response.Body())))
		console.PrintVerbose("")
		return false
	}

	// -- Console
	console.Print("\tPASSED!")

	return true
}

// GetPath Gets the path of the test file.
func (testFile *TestFile) GetPath() string {
	return filepath.Dir(testFile.filePath)
}

// GetFilePath Gets the path to the test file.
func (testFile *TestFile) GetFilePath() string {
	return testFile.filePath
}

// getRestyClient Gets a resty client.
func (testFile *TestFile) getRestyClient() *resty.Request {

	// Turn off debug
	resty.SetDebug(false)

	// Create client
	client := resty.NewRequest().
		SetContentLength(true)

	// Set headers
	for _, header := range testFile.RequestHeaders {
		client = client.SetHeader(header.Key, header.Value)
	}

	return client
}
