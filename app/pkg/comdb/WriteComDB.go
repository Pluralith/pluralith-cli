package comdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/dblock"
	"time"
)

func WriteComDB(updatedDB ComDB) error {
	functionName := "WriteComDB"

	var writeRetries int = 0
	var lock bool = true // Initializing local lock -> default: locked
	var lockObject dblock.Lock

	// Stringify updated DB for write
	updatedDBString, marshalErr := json.MarshalIndent(updatedDB, "", " ")
	if marshalErr != nil {
		return fmt.Errorf("%v: %w", functionName, marshalErr)
	}

	// Wait until DB is unlocked
	for lock && writeRetries <= 10 {
		if readErr := ReadComFile(auxiliary.PathInstance.LockPath, &lockObject); readErr != nil {
			return fmt.Errorf("reading Lock failed -> %v: %w", functionName, readErr)
		}

		// Unlocking writes for current process if lock in file belongs to current process id
		if lockObject.Id == dblock.LockInstance.Id {
			lock = false
		} else {
			// Unlocking writes for current process if comDB is not locked
			lock = lockObject.Lock
		}

		// Increment retries
		writeRetries += 1
		// Introduce random delay until next try (between 50 and 100 ms)
		time.Sleep(time.Duration(rand.Intn(100-50)+50) * time.Millisecond)
	}

	// If parse retries hit 11
	if writeRetries == 11 {
		return errors.New("writing to ComDB failed")
	}

	// Lock comDB before write
	if lockErr := dblock.UpdateDBLock(true); lockErr != nil {
		return fmt.Errorf("failed to updated lock -> %v: %w", functionName, lockErr)
	}

	// Write to Pluralith UI bus file (WriteFile replaces all file contents)
	if writeErr := os.WriteFile(auxiliary.PathInstance.ComDBPath, updatedDBString, 0700); writeErr != nil {
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	// Unlock comDB after write
	if unlockErr := dblock.UpdateDBLock(false); unlockErr != nil {
		return fmt.Errorf("failed to updated lock -> %v: %w", functionName, unlockErr)
	}

	return nil
}
