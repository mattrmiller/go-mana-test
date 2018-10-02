// Package main is the main package.
package main

// Imports
import (
	"github.com/jawher/mow.cli"
	"github.com/mattrmiller/go-mana-test/app"
	"github.com/mattrmiller/go-mana-test/console"
	"os"
)

// Constants.
const (
	AppName = "go-mana-test"
)

// Variables replaced during build.
var (
	AppVersion    string
	AppCommitHash string
	AppBuildDate  string
)

// Console
var cns = console.NewConsole()

// Main execution point
func main() {

	// Setup app
	app := cli.App(AppName, "Making APIs Nice Again - Testing Framework")
	app.Spec = "[ -c ]"

	// Define our options
	optColor := app.BoolOpt("c color", false, "Outputs console in color mode.")

	// Define our commands
	app.Command("version", "Shows version info", cmdVersion)
	app.Command("validate", "Run validation of test files", cmdValidate)
	app.Command("test", "Run tests", cmdTest)

	// Before hook
	app.Before = func() {

		// Set console color
		cns.SetOptColor(*optColor)
	}

	// Run commands
	err := app.Run(os.Args)
	if err != nil {
		cns.PrintError("Error running application")
		os.Exit(1)
	}
}

// cmdVersion Displays version info
func cmdVersion(cmd *cli.Cmd) {

	// Action
	cmd.Action = func() {

		// Setup validation app
		app := app.NewAppVersion(cns, AppVersion, AppCommitHash, AppBuildDate)
		app.Run()
	}
}

// cmdValidate Run validation
func cmdValidate(cmd *cli.Cmd) {

	// Arguments
	cmd.Spec = "PATH"
	var (
		pathProj = cmd.StringArg("PATH", "", "Path to project")
	)

	// Action
	cmd.Action = func() {

		// Setup validation app
		app := app.NewAppValidate(cns, *pathProj)
		app.Run()
	}
}

// cmdTest Run tests
func cmdTest(cmd *cli.Cmd) {

	// Arguments
	cmd.Spec = "[-bhp] PATH"
	var (
		pathProj    = cmd.StringArg("PATH", "", "Path to project")
		optBodies   = cmd.BoolOpt("b bodies", false, "Outputs HTTP request and response bodies.")
		optHault    = cmd.BoolOpt("h hault", false, "Haults on failed tests.")
		optResTimes = cmd.BoolOpt("p perf", false, "Reports HTTP response time performance.")
	)

	// Action
	cmd.Action = func() {

		// Setup validation app
		app := app.NewAppTest(cns, *pathProj, *optBodies, *optHault, *optResTimes)
		app.Run()
	}

}
