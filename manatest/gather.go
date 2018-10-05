// Package manatest provides internal workings for go-mana-test.
package manatest

// Imports
import (
	"fmt"
	"os"
	"path"
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
	files := make([]TestFile, 0)
	err := filepath.Walk(pathRead, func(pathFile string, fileInfo os.FileInfo, _ error) error {

		// Read directories and sort yaml files by name
		if fileInfo.IsDir() {

			// Yaml files
			yamlFiles, err := dirReadYamlFiles(pathFile)
			if err != nil {
				return err
			}

			// Read test files
			for _, yamlFile := range yamlFiles {

				// Read the test file
				testFile, err := ReadTestFile(path.Join(pathFile, yamlFile.Name()))
				if err != nil {
					return err
				}

				// Add
				files = append(files, *testFile)

			}

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

// dirReadYamlFiles Reads finds and sorts yaml files for a single directory.
func dirReadYamlFiles(pathDir string) ([]os.FileInfo, error) {

	// Files for return
	var ret []os.FileInfo

	// Open directory
	fileDir, err := os.Open(pathDir)
	if err != nil {
		return nil, err
	}

	// Read directory
	dirInfo, err := fileDir.Readdir(-1)
	fileDir.Close()
	if err != nil {
		return nil, err
	}

	// Iterate through each file
	for _, fileInfo := range dirInfo {

		// Only if file
		if !fileInfo.IsDir() {

			// Check if yml
			isYml, err := filepath.Match("*.yml", fileInfo.Name())
			if err != nil {
				return nil, err
			}

			// Only yml files
			if isYml {

				// Append
				ret = append(ret, fileInfo)
			}
		}
	}

	// Sort test files by name
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name() <= ret[j].Name()
	})

	return ret, nil
}
