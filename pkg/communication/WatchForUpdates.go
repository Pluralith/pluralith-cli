package communication

import (
	"os"
	"path"
	"pluralith/pkg/auxiliary"

	"github.com/fsnotify/fsnotify"
)

func WatchForUpdates() (bool, error) {
	// Set up path variables
	workingDir, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()
	pluralithDir := path.Join(homeDir, "Pluralith")
	busFilePath := path.Join(pluralithDir, "pluralith_cli.bus")

	// Create parent directories for path if they don't exist yet
	if mkErr := os.MkdirAll(pluralithDir, 0700); mkErr != nil {
		return false, mkErr
	}

	// Create file if it doesn't exist yet
	if _, fileMkErr := os.Create(busFilePath); fileMkErr != nil {
		return false, fileMkErr
	}

	// Define file watcher
	watcherInstance, watchErr := fsnotify.NewWatcher()
	if watchErr != nil {
		return false, watchErr
	}
	defer watcherInstance.Close()

	// Add bus file to watcher
	addErr := watcherInstance.Add(busFilePath)
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
				// Read bus file content
				readData, readErr := os.ReadFile(busFilePath)
				if readErr != nil {
					return false, readErr
				}

				// Parse bus file content to JSON
				parsedData, parseErr := auxiliary.ParseJson(string(readData))
				if parseErr != nil {
					return false, parseErr
				}

				// If path of latest bus file update is current working directory -> matching terraform project
				if parsedData["Path"] == workingDir {
					// If event is "confirmed" -> Execute apply, otherwise -> cancel
					if parsedData["Command"] == "confirmed" {
						return true, nil
					} else {
						return false, nil
					}
				}
			}
		// Handle watcher error
		case err := <-watcherInstance.Errors:
			return false, err
		}
	}
}
