package dblock

import (
	"fmt"
	"os"
	"path"
)

func WriteDBLock(lockString string) error {
	functionName := "WriteDBLock"

	// Generate proper path
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf("%v: %w", functionName, homeErr)
	}

	pluralithLock := path.Join(homeDir, "Pluralith", "pluralith_bus.lock")

	// Write lock string to lock file
	if writeErr := os.WriteFile(pluralithLock, []byte(lockString), 0700); writeErr != nil {
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	return nil
}
