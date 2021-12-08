package dblock

import (
	"fmt"
	"os"
	"path"
)

func InitDBLock() error {
	functionName := "InitDBLock"
	// Construct path to bus file
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf("%v: %w", functionName, homeErr)
	}

	pluralithLock := path.Join(homeDir, "Pluralith", "pluralith_bus.lock")

	// Get lock string
	lockString, getErr := LockInstance.GetLockString()
	if getErr != nil {
		return fmt.Errorf("could not get lock string -> %v: %w", functionName, getErr)
	}

	// Write to DB template to bus file
	if writeErr := os.WriteFile(pluralithLock, []byte(lockString), 0700); writeErr != nil {
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	return nil
}
