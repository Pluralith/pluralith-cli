package comdb

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"pluralith/pkg/dblock"
	"time"
)

func WriteComDB(updatedDB ComDB) error {
	var writeRetries int = 0
	var lock bool = true // Initializing local lock -> default: locked

	// Generate proper path
	homeDir, _ := os.UserHomeDir()
	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	// Stringify updated DB for write
	updatedDBString, marshalErr := json.MarshalIndent(updatedDB, "", " ")
	if marshalErr != nil {
		fmt.Println(marshalErr.Error())
		return marshalErr
	}

	// Wait until DB is unlocked
	for lock && writeRetries < 10 {
		lockObject, readErr := dblock.ReadDBLock()
		if readErr != nil {
			return readErr
		}

		// Unlocking writes for current process if lock in file belongs to current process id
		if lockObject.Id == dblock.LockInstance.Id {
			lock = false
		}

		// Unlocking writes for current process if comDB is not locked
		lock = lockObject.Lock

		// Increment retries
		writeRetries += 1
		// Introduce random delay until next try (between 50 and 100 ms)
		time.Sleep(time.Duration(rand.Intn(100-50)+50) * time.Millisecond)
	}

	// Lock comDB before write
	if lockErr := dblock.UpdateDBLock(true); lockErr != nil {
		return lockErr
	}

	// Write to Pluralith UI bus file (WriteFile replaces all file contents)
	if writeErr := os.WriteFile(pluralithBus, updatedDBString, 0700); writeErr != nil {
		return writeErr
	}

	// Unlock comDB after write
	if unlockErr := dblock.UpdateDBLock(false); unlockErr != nil {
		return unlockErr
	}

	return nil
}
