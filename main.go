// Package main is the main package.
package main

// Imports
import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/mattrmiller/go-mana-test/console"
	"github.com/mattrmiller/go-mana-test/manatest"
	"os"
	"path"
)

// Constants
const (
	APP_NAME    = "go-mana-test"
	APP_VERSION = "1.0.0"
)

// Main execution point
func main() {

	// Setup app
	app := cli.App(APP_NAME, "The Make Apis Nice Again. Testing framework.")
	verbose := app.BoolOpt("v verbose", false, "Verbose mode")

	// Handle verbose mode
	app.Before = func() {
		if *verbose {
			console.SetVerboseMode()
		}
	}

	// Define our command
	app.Command("version", "Shows version info", cmdVersion)
	app.Command("validate", "Validate tests", cmdValidate)
	app.Command("test", "Run tests", cmdTest)

	// Run commands
	app.Run(os.Args)
}

// cmdVersion Displays version info
func cmdVersion(cmd *cli.Cmd) {

	// Action
	cmd.Action = func() {

		// Print version
		console.Print(fmt.Sprintf("%s - %s\n", APP_NAME, APP_VERSION))
	}
}

// cmdValidate Validates tests
func cmdValidate(cmd *cli.Cmd) {

	// Arguments
	cmd.Spec = "PATH"
	var (
		pathProj = cmd.StringArg("PATH", "", "Path to project")
	)

	// Action
	cmd.Action = func() {

		// Read project file
		projFile, err := manatest.ReadProjectFile(*pathProj)
		if err != nil {
			console.Print(fmt.Sprintf("%s\n", err))
			os.Exit(1)
		}
		err = projFile.Validate()
		if err != nil {
			console.Print(fmt.Sprintf("Error in project file: %s\n\t%s", *pathProj, err))
		}

		// Gather test files
		pathTests := path.Join(projFile.GetPath(), projFile.Tests)
		testFiles, err := manatest.GatherTestFiles(pathTests)
		if err != nil {
			console.Print(fmt.Sprintf("%s\n", err))
			os.Exit(1)
		}

		// Validate all files
		for _, fileTest := range testFiles {

			// Validate the test file
			err = fileTest.Validate()
			if err != nil {
				console.Print(fmt.Sprintf("Error in test file: %s\n\t%s", fileTest.GetFilePath(), err))
			}
		}

		// Success
		console.Print(fmt.Sprintf("All %d test files validated!", len(testFiles)))
	}
}

// cmdTest Run tests
func cmdTest(cmd *cli.Cmd) {

	// Arguments
	cmd.Spec = "PATH"
	var (
		pathProj = cmd.StringArg("PATH", "", "Path to project")
	)

	// Action
	cmd.Action = func() {

		// Read project file
		projFile, err := manatest.ReadProjectFile(*pathProj)
		if err != nil {
			console.Print(fmt.Sprintf("%s\n", err))
			os.Exit(1)
		}
		err = projFile.Validate()
		if err != nil {
			console.Print(fmt.Sprintf("Error in project file: %s\n\t%s", *pathProj, err))
		}

		// Gather test files
		pathTests := path.Join(projFile.GetPath(), projFile.Tests)
		testFiles, err := manatest.GatherTestFiles(pathTests)
		if err != nil {
			console.Print(fmt.Sprintf("%s\n", err))
			os.Exit(1)
		}

		// Validate all files
		for _, fileTest := range testFiles {

			// Run the test file
			err = fileTest.Test(projFile)
			if err != nil {
				console.Print(fmt.Sprintf("Error running test file: %s\n\t%s", fileTest.GetFilePath(), err))
			}
		}

	}
}
