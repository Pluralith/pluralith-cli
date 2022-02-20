package dblock

import (
	"fmt"
	"os"
	"pluralith/pkg/auxiliary"
)

func UpdateDBLock(lock bool) error {
	functionName := "UpdateDBLock"
	// Set lock value and retrieve lock string
	lockString, lockErr := LockInstance.SetLock(lock)
	if lockErr != nil {
		return fmt.Errorf("could not set lock value -> %v: %w", functionName, lockErr)
	}

	// Write lock string to lock file
	if writeErr := os.WriteFile(auxiliary.StateInstance.LockPath, []byte(lockString), 0700); writeErr != nil {
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	return nil
}
