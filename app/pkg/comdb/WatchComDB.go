package comdb

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

func WatchComDB() (bool, error) {
	functionName := "WatchComDB"

	// Set up path variables
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, workingErr)
	}

	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, homeErr)
	}

	pluralithDir := path.Join(homeDir, "Pluralith")
	pluralithBus := path.Join(pluralithDir, "pluralith_bus.json")

	// Create parent directories for path if they don't exist yet
	if mkErr := os.MkdirAll(pluralithDir, 0700); mkErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, mkErr)
	}

	// Check if bus file already exists
	_, existErr := os.Stat(pluralithBus)
	if errors.Is(existErr, os.ErrNotExist) {
		// Create file if it doesn't exist yet
		if _, fileMkErr := os.Create(pluralithBus); fileMkErr != nil {
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
	addErr := watcherInstance.Add(pluralithBus)
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
				// Read comDB from file
				eventDB, readErr := ReadComDB()
				if readErr != nil {
					return false, fmt.Errorf("reading ComDB failed -> %v: %w", functionName, readErr)
				}

				// Iteratve over comDB events
				for _, event := range eventDB.Events {
					// Filter for confirm events (the only events targeted at CLI)
					if event.Path == workingDir && event.Receiver == "CLI" && !event.Received {
						// Mark event as received in comDB
						markErr := MarkComDBReceived(event)
						if markErr != nil {
							return false, fmt.Errorf("could not mark event as received -> %v --- %v: %w", event.Type, functionName, readErr)
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
