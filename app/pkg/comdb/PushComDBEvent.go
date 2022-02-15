package comdb

import (
	"fmt"
	"pluralith/pkg/auxiliary"
)

func PushComDBEvent(message ComDBEvent) error {
	functionName := "PushComDBEvent"
	var comDB ComDB

	AcquireDBLock()

	// Read DB from disk
	if readErr := ReadComFile(auxiliary.StateInstance.ComDBPath, &comDB); readErr != nil {
		return fmt.Errorf("reading ComDB failed -> %v: %w", functionName, readErr)
	}

	// Prepend new event to existing event list
	comDB.Events = append([]ComDBEvent{message}, comDB.Events...)

	// Write updated DB to disk
	if writeErr := WriteComDB(comDB); writeErr != nil {
		return fmt.Errorf("writing to ComDB failed -> %v: %w", functionName, writeErr)
	}

	return nil
}
