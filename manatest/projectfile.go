// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ProjectGlobal is a structure to handle a global variables for a project.
type ProjectGlobal struct {

	// Key, hold the key of the variable.
	Key string

	// Value, hold the value of the variable.
	Value string
}

// ProjectFile is a structure to handle a project file.
type ProjectFile struct {

	// filePath, holds the path to the file.
	filePath string

	// Name, holds the name of the test.
	Name string `yaml:"name"`

	// Tests, holds the paths to tests.
	Tests string `yaml:"tests"`

	// Globals, holds the global variables.
	Globals []ProjectGlobal `yaml:"globals"`
}

// ReadProjectFile Reads a project file.
func ReadProjectFile(pathFile string) (*ProjectFile, error) {

	// Check if yml
	if filepath.Ext(strings.TrimSpace(pathFile)) != ".yml" {
		return nil, fmt.Errorf("Project file is not a `yml` file: %s", pathFile)
	}

	// Read file
	source, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read project file at: %s", pathFile)
	}

	// Unmarshal yaml
	var projFile ProjectFile
	err = yaml.Unmarshal(source, &projFile)
	if err != nil {
		return nil, fmt.Errorf("Invalid project file format at: %s", pathFile)
	}

	// Set path
	projFile.filePath = pathFile

	return &projFile, nil
}

// Validate Validates a project file is in proper format.
func (projFile *ProjectFile) Validate() error {

	// Must have a name
	if len(projFile.Name) == 0 {
		return errors.New("roject file must have `name` field")
	}

	// Must have a tests
	if len(projFile.Tests) == 0 {
		return errors.New("project file must have `tests` field")
	}

	return nil
}

// GetPath Gets the path of the project file.
func (projFile *ProjectFile) GetPath() string {
	return filepath.Dir(projFile.filePath)
}

// GetFilePath Gets the path to the project file.
func (projFile *ProjectFile) GetFilePath() string {
	return projFile.filePath
}
