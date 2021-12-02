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
	var processId int64 = 0

	// Generate proper path
	homeDir, _ := os.UserHomeDir()
	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	// Stringify updated DB for write
	updatedDBString, marshalErr := json.MarshalIndent(updatedDB, "", " ")
	if marshalErr != nil {
		fmt.Println(marshalErr.Error())
		return marshalErr
	}

	fmt.Println(dblock.LockInstance.Lock)

	// Wait until DB is unlocked
	for dblock.LockInstance.Lock && dblock.LockInstance.Id != processId && writeRetries < 10 {
		lockObject, readErr := dblock.ReadDBLock()
		if readErr != nil {
			return readErr
		}

		// Set lock to current value in file
		dblock.LockInstance.SetLock(lockObject.Lock)
		processId = dblock.LockInstance.Id
		// Introduce random delay until next try (between 50 and 100 ms)
		time.Sleep(time.Duration(rand.Intn(100-50)+50) * time.Millisecond)
	}

	// eventDB, readErr := ReadComDB()
	// 	if readErr != nil {
	// 		fmt.Println(marshalErr.Error())
	// 		return readErr
	// 	}

	// Lock DB for write
	// ToggleLockComDB(true)

	// Write to Pluralith UI bus file (WriteFile replaces all file contents)
	if writeErr := os.WriteFile(pluralithBus, updatedDBString, 0700); writeErr != nil {
		return writeErr
	}

	// Unlock DB after write
	// ToggleLockComDB(false)

	return nil
}
