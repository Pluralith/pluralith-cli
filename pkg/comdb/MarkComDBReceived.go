package comdb

import (
	"reflect"
)

func MarkComDBReceived(event Event) error {
	// Read comDB from file
	comDB, readErr := ReadComDB()
	if readErr != nil {
		return readErr
	}

	// Iterate over events
	for index := 0; index < len(comDB.Events); index++ {
		existingEvent := &comDB.Events[index]         // Create pointer to update event value by reference
		if reflect.DeepEqual(*existingEvent, event) { // Deep compare two objects (need to dereference pointer)
			existingEvent.Received = true
		}
	}

	// Write updated comDB to file
	WriteComDB(comDB)

	return nil
}
