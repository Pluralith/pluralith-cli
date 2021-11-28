package comdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func ToggleLockComDB(lock bool) error {
	// Construct comDB path
	homeDir, _ := os.UserHomeDir()
	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	// Read comDB from file
	comDB, readErr := ReadComDB()
	if readErr != nil {
		return readErr
	}

	// Set lock value
	comDB.Locked = lock

	// Stringify updated DB for write
	comDBString, marshalErr := json.MarshalIndent(comDB, "", " ")
	if marshalErr != nil {
		fmt.Println(marshalErr.Error())
		return marshalErr
	}

	// Write to Pluralith UI bus file (WriteFile replaces all file contents)
	if writeErr := os.WriteFile(pluralithBus, comDBString, 0700); writeErr != nil {
		return writeErr
	}

	return nil
}
