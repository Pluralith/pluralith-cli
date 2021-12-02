package dblock

import (
	"os"
	"path"
)

func WriteDBLock(lockString string) error {
	// Generate proper path
	homeDir, _ := os.UserHomeDir()
	pluralithLock := path.Join(homeDir, "Pluralith", "pluralith_bus.lock")

	// Write lock string to lock file
	if writeErr := os.WriteFile(pluralithLock, []byte(lockString), 0700); writeErr != nil {
		return writeErr
	}

	return nil
}
