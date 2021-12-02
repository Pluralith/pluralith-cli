package dblock

import (
	"encoding/json"
	"os"
	"path"
)

func ReadDBLock() (Lock, error) {
	// Initialize variables
	var lockObject Lock

	// Generate proper path
	homeDir, _ := os.UserHomeDir()
	pluralithLock := path.Join(homeDir, "Pluralith", "pluralith_bus.lock")

	// Read lock file
	lockBytes, readErr := os.ReadFile(pluralithLock)
	if readErr != nil {
		// Initialize new lock if file doesn't exist
		if initErr := InitDBLock(); initErr != nil {
			return Lock{}, initErr
		}
		return Lock{}, readErr
	}

	// Unmarshal lock
	unmarshalErr := json.Unmarshal(lockBytes, &lockObject)
	if unmarshalErr != nil {
		// Initialize new lock if content of lock file is corrupted
		if initErr := InitDBLock(); initErr != nil {
			return Lock{}, initErr
		}
		return Lock{}, unmarshalErr
	}

	// Return lock string
	return lockObject, nil
}
