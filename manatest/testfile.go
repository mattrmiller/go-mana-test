// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"errors"
	"fmt"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TestHeader is a structure to handle headers for a test.
type TestHeader struct {

	// Key, hold the key of the header.
	Key string

	// Value, hold the value of the header.
	Value string
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

	// Headers, holds the header variables.
	Headers []TestHeader `yaml:"headers"`
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
	testFile.Method = strings.ToUpper(testFile.Method)
	if !ValidateMethod(testFile.Method) {
		return errors.New("Test file has invalid `method` field.")
	}

	return nil
}

// Test runs the test
func (testFile *TestFile) Test(projFile *ProjectFile) error {

	// Run for GET
	if testFile.Method == METHOD_GET {
		return testFile.runTestGet(projFile)
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

	// Create client
	client := resty.NewRequest().
		SetContentLength(true)

	// Set headers
	for _, header := range testFile.Headers {

		// Replace with globals
		if strings.Contains(header.Value, "{{") {
			header.Value = ReplaceGlobalVars(header.Value, projFile.Globals)
		}

		// Set header
		client = client.SetHeader(header.Key, header.Value)
	}

	return client
}

// Runs the test for Get tests.
func (testFile *TestFile) runTestGet(projFile *ProjectFile) error {

	// Lets for the Url, with substitutions
	url := ReplaceGlobalVars(testFile.Url, projFile.Globals)

	// Get client
	client := testFile.getRestyClient(projFile)

	// Run GET
	_, err := client.Get(url)
	if err != nil {
		return err
	}

	return nil
}
