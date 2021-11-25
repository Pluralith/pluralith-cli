package comdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
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

	// Parse DB string and handle parse error
	parseErr := json.Unmarshal([]byte(eventDBString), &eventDB)
	if parseErr != nil {
		var newErr error

		eventDB, newErr = InitComDB() // Create empty DB file
		if newErr != nil {
			fmt.Println(newErr.Error())
			return ComDB{}, newErr
		}

		return eventDB, nil
	}

	// Return parsed DB content
	return eventDB, nil
}
