package comdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"pluralith/pkg/dblock"
)

func InitComDB() (ComDB, error) {
	functionName := "InitComDB"

	// Construct path to bus file
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return ComDB{}, fmt.Errorf("%v: %w", functionName, homeErr)
	}

	pluralithBus := path.Join(homeDir, "Pluralith", "pluralith_bus.json")

	// Create empty DB template to write to file
	emptyDB := ComDB{
		Events: make([]Event, 0),
		Errors: make([]map[string]interface{}, 0),
	}

	// Turn emtpy DB template into string
	emptyDBString, marshalErr := json.MarshalIndent(emptyDB, "", " ")
	if marshalErr != nil {
		return ComDB{}, fmt.Errorf("could not format json string -> %v: %w", functionName, marshalErr)
	}

	// Write to DB template to bus file
	if writeErr := os.WriteFile(pluralithBus, emptyDBString, 0700); writeErr != nil {
		return ComDB{}, fmt.Errorf("%v: %w", functionName, writeErr)
	}

	// Instantiate comDB lock
	dblock.InitDBLock()

	return emptyDB, nil
}
