package strip

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"pluralith/pkg/ux"
)

// - - - Code to find and load files in current working directory - - -

// Function to load file content
func loadFileContent(path string) string {
	// Opening file at specified path
	openPath, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to open file")
	}

	// Reading file content
	fileByteContent, err := ioutil.ReadAll(openPath)
	if err != nil {
		log.Fatal("Failed to read file content")
	}

	return string(fileByteContent)

}

// Function to fetch files with specified extension in working directory
func FetchFiles(targetExtension string) map[string]string {
	// Instantiating new spinner
	fetchSpinner := ux.NewSpinner("Fetching State Files", "Files Fetched", "Fetching files failed")
	fetchSpinner.Start()
	// Initializing empty map to house file string content
	fileStrings := make(map[string]string)

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fetchSpinner.Fail()
		log.Fatal(err)
	}

	// Reading current working directory
	dirContent, err := ioutil.ReadDir(cwd)
	if err != nil {
		fetchSpinner.Fail()
		log.Fatal(err)
	}

	// Looping over detected files and filter for specified extension
	for _, item := range dirContent {
		file := item.Name()
		extension := filepath.Ext(file)
		name := file[0 : len(file)-len(extension)]
		// filtering for given file extension
		if filepath.Ext(file) == targetExtension {
			// Construct full path with name and working directory
			fullPath := filepath.Join(cwd, file)
			// Calling auxiliary function to load file content
			fileStrings[name] = loadFileContent(fullPath)
		}
	}

	fetchSpinner.Success(fmt.Sprintf("State Files Found: %d", len(fileStrings)))

	return fileStrings
}
