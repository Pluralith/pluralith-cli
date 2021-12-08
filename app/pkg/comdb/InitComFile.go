package comdb

import (
	"encoding/json"
	"fmt"
	"os"
	"pluralith/pkg/dblock"
	"strings"
)

func InitComFile(file string, object interface{}) error {
	functionName := "InitComFile"

	if strings.Contains(file, "ComDB") {
		object = ComDB{
			Events: make([]Event, 0),
		}
	} else {
		object = dblock.LockInstance
	}

	// Turn emtpy DB template into string
	objectString, marshalErr := json.MarshalIndent(object, "", " ")
	if marshalErr != nil {
		return fmt.Errorf("could not format json string -> %v: %w", functionName, marshalErr)
	}

	// Write to DB template to bus file
	if writeErr := os.WriteFile(file, objectString, 0700); writeErr != nil {
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	return nil
}
