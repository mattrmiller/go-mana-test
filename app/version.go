// Package app provides the command line applications commands.
package app

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-mana-test/console"
	"os"
)

// AppVersion structure handles all things related to the version app.
type AppVersion struct {

	// cns holds the console structure.
	cns *console.Console

	// appVersion holds the application version.
	appVersion string

	// appCommitHash holds the application commit hash.
	appCommitHash string

	// appBuildDate holds the application build date.
	appBuildDate string
}

// NewAppVersion creates a new version app.
func NewAppVersion(cns *console.Console, appVer string, appCmtHash string, appBuildDate string) *AppVersion {

	return &AppVersion{
		cns:           cns,
		appVersion:    appVer,
		appCommitHash: appCmtHash,
		appBuildDate:  appBuildDate,
	}
}

// Run runs the version app.
func (app *AppVersion) Run() {

	// Print version
	app.cns.PrintColor(fmt.Sprintf("Version: %s", app.appVersion), console.ColorYellow)
	app.cns.PrintColor(fmt.Sprintf("Build Date: %s", app.appBuildDate), console.ColorYellow)
	app.cns.PrintColor(fmt.Sprintf("Git Hash: %s", app.appCommitHash), console.ColorYellow)
	os.Exit(0)
}
