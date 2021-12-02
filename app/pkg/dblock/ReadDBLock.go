package dblock

import (
	"encoding/json"
	"os"
	"path"
)

func ReadDBLock() (Lock, error) {
	// Generate proper path
	homeDir, _ := os.UserHomeDir()
	pluralithLock := path.Join(homeDir, "Pluralith", "pluralith_bus.lock")

	// Read lock file
	lockBytes, readErr := os.ReadFile(pluralithLock)
	if readErr != nil {
		return Lock{}, readErr
	}

	var lockObject Lock

	// Unmarshal lock
	unmarshalErr := json.Unmarshal(lockBytes, &lockObject)
	if unmarshalErr != nil {
		return Lock{}, unmarshalErr
	}

	// Return lock string
	return lockObject, nil
}
