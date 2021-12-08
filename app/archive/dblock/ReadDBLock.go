package dblock

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func ReadDBLock() (Lock, error) {
	functionName := "ReadDBLock"

	// Initialize variables
	var lockBytes []byte
	var lockObject Lock
	var readRetries int = 0

	// Generate proper path
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return Lock{}, fmt.Errorf("%v: %w", functionName, homeErr)
	}

	pluralithLock := path.Join(homeDir, "Pluralith", "pluralith_bus.lock")

	for readRetries <= 10 {
		var readErr error

		// Read lock file
		lockBytes, readErr = os.ReadFile(pluralithLock) // TODO: Check if OS locked file
		if readErr != nil {
			// Initialize new lock if file doesn't exist
			if initErr := InitDBLock(); initErr != nil {
				return Lock{}, fmt.Errorf("initializing lock failed -> %v: %w", functionName, initErr)
			}
			return Lock{}, fmt.Errorf("%v: %w", functionName, readErr)
		}

	}

	// Unmarshal lock
	unmarshalErr := json.Unmarshal(lockBytes, &lockObject)
	if unmarshalErr != nil {
		// Initialize new lock if content of lock file is corrupted
		if initErr := InitDBLock(); initErr != nil {
			return Lock{}, fmt.Errorf("initializing ComDB failed -> %v: %w", functionName, initErr)
		}
		return Lock{}, fmt.Errorf("%v: %w", functionName, unmarshalErr)
	}

	// Return lock string
	return lockObject, nil
}
