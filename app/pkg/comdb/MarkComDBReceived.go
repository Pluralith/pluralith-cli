package comdb

import (
	"fmt"
	"reflect"
)

func MarkComDBReceived(event Event) error {
	functionName := "MarkComDBReceived"

	// Read comDB from file
	comDB, readErr := ReadComDB()
	if readErr != nil {
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
