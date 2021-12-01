package comdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

func ReadComDB() (ComDB, error) {
	// Initialize variables
	var eventDB ComDB

	// Generate proper path
	homeDir, _ := os.UserHomeDir()
	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	// Read DB file and handle non-existence
	eventDBString, readErr := os.ReadFile(pluralithBus)
	if readErr != nil {
		var newErr error

		eventDB, newErr = InitComDB() // Create empty DB file
		if newErr != nil {
			fmt.Println(newErr.Error())
			return ComDB{}, newErr
		}

		return eventDB, nil
	}

	// Initialize variables for retry logic
	// var parseSuccess bool = false
	var parseRetries int = 0

	// Retry parsing as long as criteria are met
	for parseRetries <= 10 {
		// Get file length to check if comDB file is empty
		fileSize, statErr := os.Stat(pluralithBus)
		if statErr != nil {
			return ComDB{}, statErr
		}

		// Parse DB string and handle parse error
		parseErr := json.Unmarshal([]byte(eventDBString), &eventDB)
		if parseErr != nil && fileSize.Size() == 0 { // If parsing fails and file is empty -> New comDB needs to initialized
			var newErr error

			eventDB, newErr = InitComDB() // Create empty DB file
			if newErr != nil {
				fmt.Println(newErr.Error())
				return ComDB{}, newErr
			}

			return eventDB, nil
		}

		// Increment retries
		parseRetries += 1
		// Introduce delay to avoid unnecessarily aggressive polling
		time.Sleep(50 * time.Millisecond)
	}

	// Return parsed DB content
	return eventDB, nil
}
