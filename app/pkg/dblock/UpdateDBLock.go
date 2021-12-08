package dblock

import (
	"fmt"
	"math/rand"
	"os"
	"pluralith/pkg/auxiliary"
	"time"
)

func UpdateDBLock(lock bool) error {
	functionName := "UpdateDBLock"
	// Set lock value and retrieve lock string
	lockString, lockErr := LockInstance.SetLock(lock)
	if lockErr != nil {
		return fmt.Errorf("could not set lock value -> %v: %w", functionName, lockErr)
	}

	// Write lock string to lock file
	if writeErr := os.WriteFile(auxiliary.PathInstance.LockPath, []byte(lockString), 0700); writeErr != nil {
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	// Introduce random delay until next try (between 50 and 100 ms)
	time.Sleep(time.Duration(rand.Intn(300-100)+100) * time.Millisecond)

	return nil
}
