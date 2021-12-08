package comdb

import (
	"fmt"
	"pluralith/pkg/auxiliary"
)

func PushComDBEvent(message Event) error {
	functionName := "PushComDBEvent"

	var comDB ComDB

	// Read DB from disk
	if readErr := ReadComFile(auxiliary.PathInstance.ComDBPath, &comDB); readErr != nil {
		return fmt.Errorf("reading ComDB failed -> %v: %w", functionName, readErr)
	}

	// Prepend new event to existing event list
	comDB.Events = append([]Event{message}, comDB.Events...)

	// Write updated DB to disk
	if writeErr := WriteComDB(comDB); writeErr != nil {
		return fmt.Errorf("writing to ComDB failed -> %v: %w", functionName, writeErr)
	}

	return nil
}
