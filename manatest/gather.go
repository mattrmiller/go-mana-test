// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"fmt"
	"github.com/mattrmiller/go-mana-test/console"
	"os"
	"path/filepath"
	"sort"
)

// GatherTestFiles Gathers all test files at a path.
func GatherTestFiles(pathRead string) ([]TestFile, error) {

	// Make sure directory exists
	if _, err := os.Stat(pathRead); os.IsNotExist(err) {
		return nil, fmt.Errorf("Invalid path at: %s", pathRead)
	}

	// Walk path
	files := []TestFile{}
	err := filepath.Walk(pathRead, func(pathFile string, fileInfo os.FileInfo, _ error) error {

		// Check if yml
		isYml, err := filepath.Match("*.yml", fileInfo.Name())
		if err != nil {
			return err
		}

		// Only yml files
		if isYml {

			// Read the test file
			testFile, err := ReadTestFile(pathFile)
			if err != nil {
				return err
			}

			// Add
			files = append(files, *testFile)

			// Console print
			console.PrintVerbose(fmt.Sprintf("Found test file at: %s", pathFile))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort test files by index
	sort.Slice(files, func(i, j int) bool {
		return files[i].Index <= files[j].Index
	})

	return files, nil
}
