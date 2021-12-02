package comdb

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"
)

func ReadComDB() (ComDB, error) {
	// Initialize variables
	var comDB ComDB
	var comDBRaw []byte
	var readRetries int = 0
	var parseRetries int = 0

	// Generate proper path
	homeDir, _ := os.UserHomeDir()
	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	// Read file
	// -> If read fails after 10 retries
	//		-> File does not exist
	//			-> Init ComDB
	// -> Else
	// 		-> Unmarshal file
	// 				-> If unmarshal fails after 10 retries
	//					-> File doesn't contain proper ComDB
	//						-> Init ComDB

	// Attempt to read file
	for readRetries <= 10 {
		var readErr error

		comDBRaw, readErr = os.ReadFile(pluralithBus)
		if readErr != nil {
			readRetries += 1
		} else {
			break
		}

		// Introduce delay to avoid unnecessarily aggressive polling
		time.Sleep(50 * time.Millisecond)
	}

	// If read retries hit 11
	if readRetries == 11 {
		// Init new ComDB
		newDB, newErr := InitComDB()
		if newErr != nil {
			return ComDB{}, newErr
		}

		return newDB, nil
	}

	// Attempt to parse ComDB
	for parseRetries <= 10 {
		parseErr := json.Unmarshal(comDBRaw, &comDB)
		if parseErr != nil { // If parsing fails and file is empty -> New comDB needs to initialized
			// Increment retries
			parseRetries += 1
		} else {
			break
		}

		// Introduce delay to avoid unnecessarily aggressive polling
		time.Sleep(100 * time.Millisecond)
	}

	// If parse retries hit 11
	if parseRetries == 11 {
		// Get file info
		statFile, statErr := os.Stat(pluralithBus)
		if statErr != nil {
			return ComDB{}, statErr
		}

		// If file is empty -> Init new comDB
		if statFile.Size() == 0 {
			// Init new ComDB
			newDB, newErr := InitComDB()
			if newErr != nil {
				return ComDB{}, newErr
			}

			return newDB, nil
		}

		return ComDB{}, errors.New("could not parse comDB")
	}

	// Return parsed DB content
	return comDB, nil
}
