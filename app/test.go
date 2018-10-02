// Package app provides the command line applications commands.
package app

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-mana-test/console"
	"github.com/mattrmiller/go-mana-test/manatest"
	"gopkg.in/resty.v1"
	"net/http"
	"os"
	"path"
	"time"
)

// AppTest structure handles all things related to the test app.
type AppTest struct {

	// cns holds the console structure.
	cns *console.Console

	// pathProj holds the project path.
	pathProj string

	// optBodies holds the options to output HTTP bodies.
	optBodies bool

	// optHault holds the options to hault on failed tests.
	optHault bool

	// optResTimes holds the options to calculate response times.
	optResTimes bool
}

// NewAppTest creates a new test app.
func NewAppTest(cns *console.Console, pathProj string, optBodies bool, optHault bool, optResTimes bool) *AppTest {

	return &AppTest{
		cns:         cns,
		pathProj:    pathProj,
		optBodies:   optBodies,
		optHault:    optHault,
		optResTimes: optResTimes,
	}
}

// Run runs the validation app.
func (app *AppTest) Run() {

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

		// Validate the test file
		err = fileTest.Validate()
		if err != nil {
			app.cns.PrintError(fmt.Sprintf("Error in test file: %s\n\t%s", fileTest.GetFilePath(), err))
			os.Exit(1)
		}

		// Make the test headers
		testHeaders := fileTest.MakeTestHeaders(projFile)

		// Make the test URL
		testURL := fileTest.MakeTestURL(projFile)

		// Console
		app.cns.Print(fmt.Sprintf("\nRunning test: %s...", fileTest.Name))
		app.cns.Print(fmt.Sprintf("\t%s: %s", fileTest.RequestMethod, testURL))

		// Prepare resty client
		resty.SetDebug(false)
		client := resty.NewRequest().
			SetContentLength(true)
		for _, header := range testHeaders {
			client = client.SetHeader(header.Key, header.Value)
		}

		// Set body
		if fileTest.ReqBody != nil && fileTest.RequestMethod != http.MethodTrace {
			client.SetBody(manatest.ReplaceVars(fileTest.ReqBody.(string), &projFile.Globals))
		}

		// Always calculate response times to save on logic check costs
		timeStart := time.Now()

		// Run
		response, err := client.Execute(fileTest.RequestMethod, testURL)
		if err != nil {
			app.cns.PrintError(fmt.Sprintf("\tError executing HTTP request: %s", err))
			os.Exit(1)
		}

		// Calculate respjse time
		timeEnd := time.Now()
		resTimeMs := timeEnd.Sub(timeStart).Seconds() * 1000

		// Output response times?
		if app.optResTimes {
			app.cns.PrintColor(fmt.Sprintf("\n\tResponse Time: %.2fms", resTimeMs), console.ColorYellow)
		}

		// Verbose bodies
		if app.optBodies {
			if client.Body != nil {
				app.cns.PrintColor("\n\tRequest Body:", console.ColorCyan)
				app.cns.PrintColor(fmt.Sprintf("\t\t%s\n\n", client.Body.(string)), console.ColorCyan)
			}
			if response.RawResponse != nil {
				app.cns.PrintColor("\n\tResponse Body:", console.ColorCyan)
				app.cns.PrintColor(fmt.Sprintf("\t\t%s\n\n", string(response.Body())), console.ColorCyan)
			}
		}

		// Save cache
		err = manatest.SaveCacheFromResponse(&fileTest.Cache, response)
		if err != nil {
			app.cns.PrintError(fmt.Sprintf("\n\tError saving cache: %s", err))
			os.Exit(1)
		}

		// Run tests
		app.cns.Print("\tRunning checks...")
		err = manatest.RunChecks(&fileTest.Checks, &projFile.Globals, response)
		if err != nil {
			app.cns.PrintError(fmt.Sprintf("\tFAIL: %s\n", err))
			if app.optHault {
				os.Exit(1)
			}
		} else {
			app.cns.PrintColor("\tPASSED!", console.ColorGreen)
		}
	}

	// Dry run results
	app.cns.PrintColor(fmt.Sprintf("\nAll %d tests passed tests!", len(testFiles)), console.ColorGreen)
}
