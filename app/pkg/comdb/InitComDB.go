package comdb

import (
	"encoding/json"
	"os"
	"path"
	"pluralith/pkg/dblock"
)

func InitComDB() (ComDB, error) {
	// Construct path to bus file
	homeDir, _ := os.UserHomeDir()
	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	// Create empty DB template to write to file
	emptyDB := ComDB{
		Events: make([]Event, 0),
		Errors: make([]map[string]interface{}, 0),
	}

	// Turn emtpy DB template into string
	emptyDBString, marshalErr := json.MarshalIndent(emptyDB, "", " ")
	if marshalErr != nil {
		return ComDB{}, marshalErr
	}

	// Write to DB template to bus file
	if writeErr := os.WriteFile(pluralithBus, emptyDBString, 0700); writeErr != nil {
		return ComDB{}, writeErr
	}

	// Instantiate comDB lock
	dblock.InitDBLock()

	return emptyDB, nil
}
