// Package app provides the command line applications commands.
package app

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-mana-test/console"
	"github.com/mattrmiller/go-mana-test/manatest"
	"os"
	"path"
)

// AppValidate structure handles all things related to the validation app.
type AppValidate struct {

	// cns holds the console structure.
	cns *console.Console

	// pathProj holds the project path.
	pathProj string
}

// AppValidate creates a new validation app.
func NewAppValidate(cns *console.Console, pathProj string) *AppValidate {

	return &AppValidate{
		cns:      cns,
		pathProj: pathProj,
	}
}

// Run runs the validation app.
func (app *AppValidate) Run() {

	// Read project file
	projFile, err := manatest.ReadProjectFile(app.pathProj)
	if err != nil {
		app.cns.PrintError(fmt.Sprintf("Error reading project file: %s\n\t%s", app.pathProj, err))
		os.Exit(1)
	}
	err = projFile.Validate()
	if err != nil {
		app.cns.PrintError(fmt.Sprintf("Error in project file: %s\n\t%s", app.pathProj, err))
		os.Exit(1)
	}

	// Gather test files
	pathTests := path.Join(projFile.GetPath(), projFile.Tests)
	testFiles, err := manatest.GatherTestFiles(pathTests)
	if err != nil {
		app.cns.PrintError(fmt.Sprintf("Error gathering test files: %s\n", err))
		os.Exit(1)
	}

	// Validate all files
	for _, fileTest := range testFiles {

		// Print details
		app.cns.Print(fmt.Sprintf("Test File: %s", fileTest.GetFilePath()))

		// Validate the test file
		err = fileTest.Validate()
		if err != nil {
			app.cns.PrintError(fmt.Sprintf("\tValidation Result: FAIL: %s", err))
			os.Exit(1)
		}
		app.cns.PrintColor("\tValidation Result: PASSED!", console.ColorGreen)
	}

	// Results
	app.cns.PrintColor(fmt.Sprintf("\nAll %d tests passed validation!", len(testFiles)), console.ColorGreen)
	os.Exit(0)
}
