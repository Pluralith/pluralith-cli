package communication

import (
	"encoding/json"
	"os"
	"path"
	"pluralith/pkg/auxiliary"
)

func createEmptyDB(pluralithBus string) (CommunicationDB, error) {
	// Create empty db template to write to file
	emptyDB := CommunicationDB{
		Locked: false,
		Events: make([]interface{}, 0),
		Errors: make([]interface{}, 0),
	}

	// Turning emtpy db template into string
	emptyDBString, marshalErr := json.MarshalIndent(emptyDB, "", " ")
	if marshalErr != nil {
		return CommunicationDB{}, marshalErr
	}

	// Write to db template to bus file
	if writeErr := os.WriteFile(pluralithBus, emptyDBString, 0700); writeErr != nil {
		return CommunicationDB{}, writeErr
	}

	return emptyDB, nil
}

func readCommunicationDB(pluralithBus string) (CommunicationDB, bool, error) {
	// Initialize eventDB variable
	var eventDB CommunicationDB

	// Read bus file content
	readData, readErr := os.ReadFile(pluralithBus)
	if readErr != nil { // If file doesn't exist yet
		newDB, newErr := createEmptyDB(pluralithBus) // Create empty DB file
		if newErr != nil {
			return CommunicationDB{}, false, newErr
		}

		eventDB = newDB
	} else { // If file exists
		// Parse bus file content to JSON
		parsedData, parseErr := auxiliary.ParseJson(string(readData))
		if parseErr != nil { // If parsing failed (e.g. empty file)
			newDB, newErr := createEmptyDB(pluralithBus) // Create empty DB file
			if newErr != nil {
				return CommunicationDB{}, false, newErr
			}

			eventDB = newDB
		} else {
			// Map values from parsed data onto eventDB interface
			eventDB = CommunicationDB{
				Locked: parsedData["Locked"].(bool),
				Events: parsedData["Events"].([]interface{}),
				Errors: parsedData["Errors"].([]interface{}),
			}
		}
	}

	// Return lock true if spinlock active
	if eventDB.Locked {
		return eventDB, true, nil
	}

	// Return parsed DB content
	return eventDB, false, nil
}

func EmitUpdate(message Update) error {
	// Setting up directory path variables
	homeDir, _ := os.UserHomeDir()
	pluralithDir := path.Join(homeDir, "Pluralith")
	pluralithBus := path.Join(pluralithDir, "pluralith_bus.json")

	// Create parent directories for path if they don't exist yet
	if mkErr := os.MkdirAll(pluralithDir, 0700); mkErr != nil {
		return mkErr
	}

	// Initialize spinlock variable
	var eventDB CommunicationDB
	locked := true

	for locked {
		// Reading communication database
		db, lockStatus, dbErr := readCommunicationDB(pluralithBus)
		if dbErr != nil {
			return dbErr
		}

		// Retrieving lock status and database
		locked = lockStatus
		eventDB = db
	}

	// Append new event to list
	eventDB.Events = append(eventDB.Events, message)

	// JSONify DB
	dbJson, dbJsonErr := json.MarshalIndent(eventDB, "", " ")
	if dbJsonErr != nil {
		return dbJsonErr
	}

	// Write to Pluralith UI bus file (WriteFile replaces all file contents)
	if writeErr := os.WriteFile(pluralithBus, dbJson, 0700); writeErr != nil {
		return writeErr
	}

	return nil
}
