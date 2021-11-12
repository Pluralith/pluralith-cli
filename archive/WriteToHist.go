package communication

import (
	"log"
	"os"
	"path"
)

// - - - Code to write to command history - - -

func WriteToHist(command string, text string) {
	// Setting up directory path variables
	homeDir, _ := os.UserHomeDir()
	workingDir, _ := os.Getwd()
	pluralithDir := path.Join(homeDir, "Pluralith")

	var histEntry string

	if text == "" {
		histEntry = workingDir + "-----" + command + "\n"
	} else {
		histEntry = text
	}

	// Creating parent directories for path if they don't exist yet
	if err := os.MkdirAll(pluralithDir, 0700); err != nil {
		log.Fatal(err)
	}

	// Appending to file or creating it if it doesn't exist yet
	file, err := os.OpenFile(path.Join(pluralithDir, "pluralith.hist"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}

	// Appending line to file
	if _, err := file.Write([]byte(histEntry)); err != nil {
		log.Fatal(err)
	}

	// Closing file
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}
