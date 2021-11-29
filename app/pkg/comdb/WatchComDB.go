package comdb

import (
	"errors"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

func handleComDBEvent() {

}

func WatchComDB() (bool, error) {
	// Set up path variables
	workingDir, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()
	pluralithDir := path.Join(homeDir, "Pluralith")
	pluralithBus := path.Join(pluralithDir, "pluralith_bus.json")

	// Create parent directories for path if they don't exist yet
	if mkErr := os.MkdirAll(pluralithDir, 0700); mkErr != nil {
		return false, mkErr
	}

	// Check if bus file already exists
	_, existErr := os.Stat(pluralithBus)
	if errors.Is(existErr, os.ErrNotExist) {
		// Create file if it doesn't exist yet
		if _, fileMkErr := os.Create(pluralithBus); fileMkErr != nil {
			return false, fileMkErr
		}
	}

	// Define file watcher
	watcherInstance, watchErr := fsnotify.NewWatcher()
	if watchErr != nil {
		return false, watchErr
	}
	defer watcherInstance.Close()

	// Add bus file to watcher
	addErr := watcherInstance.Add(pluralithBus)
	if addErr != nil {
		return false, addErr
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
					return false, readErr
				}

				// Iteratve over comDB events
				for _, event := range eventDB.Events {
					// Filter for confirm events (the only events targeted at CLI)
					if event.Path == workingDir && event.Receiver == "CLI" && !event.Received && event.Type == "confirmed" {
						MarkComDBReceived(event) // Mark event as received in comDB
						return true, nil
					}
				}
			}
		// Handle watcher error
		case err := <-watcherInstance.Errors:
			return false, err
		}
	}
}
