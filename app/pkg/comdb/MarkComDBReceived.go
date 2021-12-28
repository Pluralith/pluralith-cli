package comdb

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"reflect"
)

func MarkComDBReceived(event ComDBEvent) error {
	functionName := "MarkComDBReceived"

	var comDB ComDB

	// Read comDB from file
	if readErr := ReadComFile(auxiliary.PathInstance.ComDBPath, &comDB); readErr != nil {
		return fmt.Errorf("reading ComDB failed -> %v: %w", functionName, readErr)
	}

	// Iterate over events
	for index := 0; index < len(comDB.Events); index++ {
		existingEvent := &comDB.Events[index]         // Create pointer to update event value by reference
		if reflect.DeepEqual(*existingEvent, event) { // Deep compare two objects (need to dereference pointer)
			existingEvent.Received = true
		}
	}

	// Write updated DB to disk
	writeErr := WriteComDB(comDB)
	if writeErr != nil {
		return fmt.Errorf("writing to ComDB failed -> %v: %w", functionName, writeErr)
	}

	return nil
}
