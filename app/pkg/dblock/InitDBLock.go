package dblock

import (
	"os"
	"path"
)

func InitDBLock() error {
	// Construct path to bus file
	homeDir, _ := os.UserHomeDir()
	pluralithLock := path.Join(homeDir, "Pluralith", "pluralith_bus.lock")

	// Get lock string
	lockString, getErr := LockInstance.GetLockString()
	if getErr != nil {
		return getErr
	}

	// Write to DB template to bus file
	if writeErr := os.WriteFile(pluralithLock, []byte(lockString), 0700); writeErr != nil {
		return writeErr
	}

	return nil
}

// func ToggleLockComDB(lock bool) error {
// 	// Construct comDB path
// 	homeDir, _ := os.UserHomeDir()
// 	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

// 	// Read comDB from file
// 	comDB, readErr := ReadComDB()
// 	if readErr != nil {
// 		return readErr
// 	}

// 	// Set lock value
// 	comDB.Locked = lock

// 	// Stringify updated DB for write
// 	comDBString, marshalErr := json.MarshalIndent(comDB, "", " ")
// 	if marshalErr != nil {
// 		fmt.Println(marshalErr.Error())
// 		return marshalErr
// 	}

// 	// Write to Pluralith UI bus file (WriteFile replaces all file contents)
// 	if writeErr := os.WriteFile(pluralithBus, comDBString, 0700); writeErr != nil {
// 		return writeErr
// 	}

// 	return nil
// }
