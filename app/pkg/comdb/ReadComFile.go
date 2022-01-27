package comdb

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func ReadComFile(file string, object interface{}) error {
	functionName := "ReadComFile"

	// Initialize variables
	var byteResult []byte
	var readRetries int = 0
	var parseRetries int = 0

	// Attempt to read file
	for readRetries <= 10 {
		var readErr error

		byteResult, readErr = os.ReadFile(file)
		if readErr != nil {
			readRetries += 1
		} else {
			break
		}

		// Introduce random delay until next try (between 50 and 100 ms)
		time.Sleep(time.Duration(rand.Intn(100-50)+50) * time.Millisecond)
	}

	// If read retries hit 11
	if readRetries == 11 {
		// Init new ComDB
		if newErr := InitComFile(file, &object); newErr != nil {
			return fmt.Errorf("initializing %v failed -> %v: %w", file, functionName, newErr)
		}

		return nil
	}

	// Attempt to parse ComDB
	for parseRetries <= 10 {
		parseErr := json.Unmarshal(byteResult, &object)
		if parseErr != nil { // If parsing fails and file is empty -> New comDB needs to initialized
			// Increment retries
			parseRetries += 1
		} else {
			break
		}

		// Introduce random delay until next try (between 50 and 100 ms)
		time.Sleep(time.Duration(rand.Intn(100-50)+50) * time.Millisecond)
	}

	// If parse retries hit 11
	if parseRetries == 11 {
		// Get file info
		statFile, statErr := os.Stat(file)
		if statErr != nil {
			return fmt.Errorf("could not get file information -> %v: %w", functionName, statErr)
		}

		// If file is empty -> Init new comDB
		if statFile.Size() == 0 {
			// Init new ComDB
			if newErr := InitComFile(file, &object); newErr != nil {
				return fmt.Errorf("initializing %v failed -> %v: %w", file, functionName, newErr)
			}

			return nil
		}

		return fmt.Errorf("initializing %v failed -> %v", file, functionName)
	}

	// Return parsed DB content
	return nil
}
