package comdb

import (
	"fmt"
)

func PushComDBEvent(message Event) error {
	// Read DB from disk
	eventDB, readErr := ReadComDB()
	if readErr != nil {
		fmt.Println(readErr.Error())
		return readErr
	}

	// Append new event to existing event list
	eventDB.Events = append(eventDB.Events, message)

	// Write updated DB to disk
	writeErr := WriteComDB(eventDB)
	if writeErr != nil {
		fmt.Println(writeErr.Error())
		return writeErr
	}

	return nil
}
