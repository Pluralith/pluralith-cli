package dblock

import (
	"math/rand"
	"time"
)

func UpdateDBLock(lock bool) error {
	// Set lock value and retrieve lock string
	lockString, lockErr := LockInstance.SetLock(lock)
	if lockErr != nil {
		return lockErr
	}

	// Write lock string to lock file
	writeErr := WriteDBLock(lockString)
	if writeErr != nil {
		return writeErr
	}

	// Introduce random delay until next try (between 50 and 100 ms)
	time.Sleep(time.Duration(rand.Intn(300-100)+100) * time.Millisecond)

	return nil
}
