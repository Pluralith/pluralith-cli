package comdb

import (
	"fmt"
	"os"
	"path"
	"pluralith/pkg/auxiliary"
)

func ReadComDB() (ComDB, error) {
	// Initialize variables
	var eventDB ComDB
	var lock bool = true

	// Generate proper path
	homeDir, _ := os.UserHomeDir()
	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	for lock {
		// Read DB file and handle non-existence
		eventDBString, readErr := os.ReadFile(pluralithBus)
		if readErr != nil {
			var newErr error

			eventDB, newErr = InitComDB() // Create empty DB file
			if newErr != nil {
				fmt.Println(newErr.Error())
				return ComDB{}, newErr
			}
		}

		// Parse DB string and handle parse error
		eventDBObject, parseErr := auxiliary.ParseJson(string(eventDBString))
		if parseErr != nil {
			var newErr error

			eventDB, newErr = InitComDB() // Create empty DB file
			if newErr != nil {
				fmt.Println(newErr.Error())
				return ComDB{}, newErr
			}
		}

		// Construct ComDB object
		eventDB = ComDB{
			Locked: eventDBObject["Locked"].(bool),
			Events: eventDBObject["Events"].([]interface{}),
			Errors: eventDBObject["Errors"].([]interface{}),
		}

		// Update lock variable
		lock = eventDBObject["Locked"].(bool)
	}

	// Return parsed DB content
	return eventDB, nil
}
