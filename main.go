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
	AppName    = "go-mana-test"
	AppVersion = "0.1.2"
)

// Main execution point
func main() {

	// Setup app
	app := cli.App(AppName, "Making APIs Nice Again - Testing Framework")
	verbose := app.BoolOpt("v verbose", false, "Verbose mode")

	// Handle verbose mode
	app.Before = func() {
		if *verbose {
			console.SetVerboseMode()
		}
	}

	// Define our command
	app.Command("version", "Shows version info", cmdVersion)
	app.Command("test", "Run tests", cmdTest)

	// Run commands
	err := app.Run(os.Args)
	if err != nil {
		console.Print("Error running application")
		os.Exit(1)
	}
}

// cmdVersion Displays version info
func cmdVersion(cmd *cli.Cmd) {

	// Action
	cmd.Action = func() {

		// Print version
		console.Print(fmt.Sprintf("%s - %s\n", AppName, AppVersion))
	}
}

// cmdTest Run tests
func cmdTest(cmd *cli.Cmd) {

	// Arguments
	cmd.Spec = "PATH [ -d ]"
	var (
		pathProj = cmd.StringArg("PATH", "", "Path to project")
		dryRun   = cmd.BoolOpt("d dryrun", false, "Dry run, only validates test files.")
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
				os.Exit(1)
			}

			// Run the test file
			if !*dryRun {
				if !fileTest.Test(projFile) {
					os.Exit(1)
				}
			}
		}

		// Dry run results
		if !*dryRun {
			console.Print(fmt.Sprintf("\nAll %d tests passed tests!", len(testFiles)))
		} else {
			console.Print(fmt.Sprintf("%d tests found and passed validation!", len(testFiles)))
		}
	}

}
