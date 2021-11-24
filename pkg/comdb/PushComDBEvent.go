package comdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func PushComDBEvent(message Update) error {
	// Initialize variables
	var eventDB ComDB
	var lock bool = true
	var readErr error

	// 1) Read DB + wait until unlocked if locked
	for lock {
		eventDB, lock, readErr = ReadComDB()
		if readErr != nil {
			fmt.Println(readErr.Error())
			return readErr
		}
	}

	// 2) Append new event
	eventDB.Events = append(eventDB.Events, message)
	// 3) WriteComDB() -> Locking handled inside Write function

	// Setting up directory path variables
	homeDir, _ := os.UserHomeDir()
	pluralithDir := path.Join(homeDir, "Pluralith")
	pluralithBus := path.Join(pluralithDir, "pluralith_bus.json")

	// Create parent directories for path if they don't exist yet
	// if mkErr := os.MkdirAll(pluralithDir, 0700); mkErr != nil {
	// 	fmt.Println(mkErr.Error())
	// 	return mkErr
	// }

	// Initialize spinlock variable
	// var eventDB ComDB
	// locked := true

	// for locked {
	// 	// Reading communication database
	// 	// db, lockStatus, dbErr := ReadComDB(pluralithBus)
	// 	// if dbErr != nil {
	// 	// 	fmt.Println(dbErr.Error())
	// 	// 	return dbErr
	// 	// }

	// 	// Retrieving lock status and database
	// 	locked = lockStatus
	// 	eventDB = db
	// }

	// Append new event to list
	eventDB.Events = append(eventDB.Events, message)

	// JSONify DB
	dbJson, dbJsonErr := json.MarshalIndent(eventDB, "", " ")
	if dbJsonErr != nil {
		fmt.Println(dbJsonErr.Error())
		return dbJsonErr
	}

	// Write to Pluralith UI bus file (WriteFile replaces all file contents)
	if writeErr := os.WriteFile(pluralithBus, dbJson, 0700); writeErr != nil {
		return writeErr
	}

	return nil
}
