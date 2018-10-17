// Package app provides the command line applications commands.
package app

// Imports
import (
	"bytes"
	"encoding/json"
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

	// cns stores the console structure.
	cns *console.Console

	// pathProj stores the project path.
	pathProj string

	// optBodies stores the options to output HTTP bodies.
	optBodies bool

	// optExit stores the options to hault on failed tests.
	optExit bool

	// optResTimes stores the options to calculate response times.
	optResTimes bool
}

// NewAppTest creates a new test app.
func NewAppTest(cns *console.Console, pathProj string, optBodies bool, optExit bool, optResTimes bool) *AppTest {

	return &AppTest{
		cns:         cns,
		pathProj:    pathProj,
		optBodies:   optBodies,
		optExit:     optExit,
		optResTimes: optResTimes,
	}
}

// Run runs the validation app.
func (app *AppTest) Run() {

	// Read project file
	projFile, err := manatest.ReadProjectFile(app.pathProj)
	if err != nil {
		app.cns.PrintError(fmt.Sprintf("Error reading project file: %s\n%s", app.pathProj, err))
		os.Exit(1)
	}
	err = projFile.Validate()
	if err != nil {
		app.cns.PrintError(fmt.Sprintf("Error in project file: %s\n%s", app.pathProj, err))
		os.Exit(1)
	}

	// Gather test files
	pathTests := path.Join(projFile.GetPath(), projFile.Tests)
	testFiles, err := manatest.GatherTestFiles(pathTests)
	if err != nil {
		app.cns.PrintError(fmt.Sprintf("Error gathering test files: %s\n", err))
		os.Exit(1)
	}

	// Response times
	var resTimeTally, resTimeMax, resTimeMin float64
	var resTimeMaxName, resTimeMinName string
	resTimeTally = 0
	resTimeMax = 0
	resTimeMin = -1

	// Test count
	var countTestSucc, countTestFail int
	countTestSucc = 0
	countTestFail = 0

	// Validate all files
	for _, fileTest := range testFiles {

		// Validate the test file
		err = fileTest.Validate()
		if err != nil {
			app.cns.PrintError(fmt.Sprintf("Error in test file: %s\n%s", fileTest.GetFilePath(), err))
			os.Exit(1)
		}

		// Make the test headers
		testHeaders := fileTest.MakeTestHeaders(projFile)

		// Make the test URL
		testURL := fileTest.MakeTestURL(projFile)

		// Console
		app.cns.Print(fmt.Sprintf("--\nRunning Test: %s...", fileTest.Name))
		app.cns.PrintColor(fmt.Sprintf("%s: %s", fileTest.RequestMethod, testURL), console.ColorMagenta)

		// Prepare resty client
		resty.SetDebug(false)
		client := resty.NewRequest().
			SetContentLength(true)
		for _, header := range testHeaders {
			client = client.SetHeader(header.Key, header.Value)
		}

		// Set body
		if fileTest.ReqBody != nil && fileTest.RequestMethod != http.MethodTrace {
			client.SetBody(manatest.ReplaceVarsInRequestBody(fileTest.ReqBody.(string), &projFile.Globals))
		}

		// Always calculate response times to save on logic check costs
		timeStart := time.Now()

		// Run
		response, err := client.Execute(fileTest.RequestMethod, testURL)
		if err != nil {
			app.cns.PrintError(fmt.Sprintf("Error executing HTTP request: %s", err))
			os.Exit(1)
		}

		// Calculate response time
		timeEnd := time.Now()
		resTimeMs := timeEnd.Sub(timeStart).Seconds() * 1000
		resTimeTally += resTimeMs
		if resTimeMs > resTimeMax {
			resTimeMax = resTimeMs
			resTimeMaxName = fileTest.Name
		}
		if resTimeMin == -1 || resTimeMs < resTimeMin {
			resTimeMin = resTimeMs
			resTimeMinName = fileTest.Name
		}

		// Output response times?
		if app.optResTimes {
			app.cns.PrintColor(fmt.Sprintf("Response Time: %.2fms", resTimeMs), console.ColorYellow)
		}

		// Verbose bodies
		if app.optBodies {
			if client.Body != nil {
				app.cns.PrintColor("Request Body:", console.ColorCyan)
				var jsonPretty bytes.Buffer
				err := json.Indent(&jsonPretty, []byte(client.Body.(string)), "", "    ")
				if err != nil {
					app.cns.PrintError(fmt.Sprintf("Error outputting request body: %s", err))
					os.Exit(1)
				}
				app.cns.PrintColor(fmt.Sprintf("%s", jsonPretty.String()), console.ColorCyan)
			}
			if response.RawResponse != nil && len(response.Body()) != 0 {
				var jsonPretty bytes.Buffer
				err := json.Indent(&jsonPretty, response.Body(), "", "    ")
				if err != nil {
					app.cns.PrintError(fmt.Sprintf("Error outputting response body: %s", err))
					os.Exit(1)
				}
				app.cns.PrintColor("Response Body:", console.ColorCyan)
				app.cns.PrintColor(fmt.Sprintf("%s", jsonPretty.String()), console.ColorCyan)
			}
		}

		// Save cache
		err = manatest.SaveCacheFromResponse(&fileTest.Cache, response)
		if err != nil {
			app.cns.PrintError(fmt.Sprintf("Error saving cache: %s", err))
			os.Exit(1)
		}

		// Run tests
		err = manatest.RunChecks(&fileTest.Checks, &projFile.Globals, response)
		if err != nil {
			countTestFail++
			app.cns.PrintError(fmt.Sprintf("Test Result: FAIL: %s", err))
			if app.optExit {
				os.Exit(1)
			}
		} else {
			countTestSucc++
			app.cns.PrintColor("Test Result: PASSED!", console.ColorGreen)
		}
	}

	// Total files
	totalFiles := len(testFiles)

	// Output response times?
	if app.optResTimes {
		app.cns.Print("--")
		app.cns.PrintColor(fmt.Sprintf("Average Response Time: %.2fms", resTimeTally/float64(totalFiles)), console.ColorYellow)
		app.cns.PrintColor(fmt.Sprintf("Response Time Max : %.2fms - %s", resTimeMax, resTimeMaxName), console.ColorYellow)
		app.cns.PrintColor(fmt.Sprintf("Response Time Min : %.2fms - %s", resTimeMin, resTimeMinName), console.ColorYellow)
	}

	// Results
	app.cns.Print("--")
	app.cns.PrintColor(fmt.Sprintf("Passing Tests: %d", countTestSucc), console.ColorGreen)
	if countTestFail != 0 {
		app.cns.PrintError(fmt.Sprintf("Failing Tests: %d", countTestFail))
		os.Exit(1)
	} else {
		os.Exit(0)
	}

}
