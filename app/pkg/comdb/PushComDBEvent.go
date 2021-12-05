package comdb

import (
	"fmt"
	"time"
)

func PushComDBEvent(message Event) error {
	functionName := "PushComDBEvent"

	// Read DB from disk
	comDB, readErr := ReadComDB()
	if readErr != nil {
		return fmt.Errorf("reading ComDB failed -> %v: %w", functionName, readErr)
	}

	// Prepend new event to existing event list
	comDB.Events = append([]Event{message}, comDB.Events...)

	// Write updated DB to disk
	writeErr := WriteComDB(comDB)
	if writeErr != nil {
		return fmt.Errorf("writing to ComDB failed -> %v: %w", functionName, writeErr)
	}

	// Introduce small delay to ensure proper event pushes
	time.Sleep(50 * time.Millisecond)

	return nil
}
