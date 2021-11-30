package comdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func WriteComDB(updatedDB ComDB) error {
	// Initialize variables
	var lock bool = true

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
	for lock {
		eventDB, readErr := ReadComDB()
		if readErr != nil {
			fmt.Println(marshalErr.Error())
			return readErr
		}

		lock = eventDB.Locked
	}

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
