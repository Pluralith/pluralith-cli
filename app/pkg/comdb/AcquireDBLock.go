package comdb

import (
	"fmt"
	"math/rand"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/dblock"
	"time"
)

func AcquireDBLock() error {
	functionName := "AcquireLock"

	var writeRetries int = 0
	var lock bool = true // Initializing local lock -> default: locked

	// Wait until DB is unlocked
	for lock && writeRetries <= 10 {
		var lockObject dblock.Lock

		if readErr := ReadComFile(auxiliary.StateInstance.LockPath, &lockObject); readErr != nil {
			return fmt.Errorf("reading Lock failed -> %v: %w", functionName, readErr)
		}

		lock = lockObject.Lock

		if lock {
			// Increment retries
			writeRetries += 1
			// Introduce random delay until next try (between 50 and 100 ms)
			time.Sleep(time.Duration(rand.Intn(100-50)+50) * time.Millisecond)
		}
	}

	// Acquire lock
	if lockErr := dblock.UpdateDBLock(true); lockErr != nil {
		return fmt.Errorf("failed to updated lock -> %v: %w", functionName, lockErr)
	}

	return nil
}
