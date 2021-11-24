package comdb

import (
	"encoding/json"
	"os"
	"path"
)

func InitComDB() (ComDB, error) {
	// Construct path to bus file
	homeDir, _ := os.UserHomeDir()
	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	// Create empty db template to write to file
	emptyDB := ComDB{
		Locked: false,
		Events: make([]interface{}, 0),
		Errors: make([]interface{}, 0),
	}

	// Turning emtpy db template into string
	emptyDBString, marshalErr := json.MarshalIndent(emptyDB, "", " ")
	if marshalErr != nil {
		return ComDB{}, marshalErr
	}

	// Write to db template to bus file
	if writeErr := os.WriteFile(pluralithBus, emptyDBString, 0700); writeErr != nil {
		return ComDB{}, writeErr
	}

	return emptyDB, nil
}
