package comdb

import (
	"errors"
	"fmt"
	"os"
	"pluralith/pkg/auxiliary"

	"github.com/fsnotify/fsnotify"
)

func WatchComDB() (bool, error) {
	functionName := "WatchComDB"

	// Check if bus file already exists
	_, existErr := os.Stat(auxiliary.PathInstance.ComDBPath)
	if errors.Is(existErr, os.ErrNotExist) {
		// Create file if it doesn't exist yet
		if _, fileMkErr := os.Create(auxiliary.PathInstance.ComDBPath); fileMkErr != nil {
			return false, fmt.Errorf("%v: %w", functionName, fileMkErr)
		}
	}

	// Define file watcher
	watcherInstance, watchErr := fsnotify.NewWatcher()
	if watchErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, watchErr)
	}
	defer watcherInstance.Close()

	// Add bus file to watcher
	addErr := watcherInstance.Add(auxiliary.PathInstance.ComDBPath)
	if addErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, addErr)
	}

	// Handle watcher events
	for {
		select {
		case event := <-watcherInstance.Events:
			// Switch different watcher event types
			switch {
			// If a write event happens
			case event.Op&fsnotify.Write == fsnotify.Write:

				var comDB ComDB
				// Read comDB from file
				if readErr := ReadComFile(auxiliary.PathInstance.ComDBPath, &comDB); readErr != nil {
					return false, fmt.Errorf("reading ComDB failed -> %v: %w", functionName, readErr)
				}

				// Iteratve over comDB events
				for _, event := range comDB.Events {
					// Filter for confirm events (the only events targeted at CLI)
					if event.Path == auxiliary.PathInstance.WorkingPath && event.Receiver == "CLI" && !event.Received {
						// Mark event as received in comDB
						if markErr := MarkComDBReceived(event); markErr != nil {
							return false, fmt.Errorf("could not mark event as received -> %v --- %v: %w", event.Type, functionName, markErr)
						}

						if event.Type == "confirmed" {
							return true, nil
						} else {
							return false, nil
						}
					}
				}
			}
		// Handle watcher error
		case err := <-watcherInstance.Errors:
			return false, fmt.Errorf("%v: %w", functionName, err)
		}
	}
}
