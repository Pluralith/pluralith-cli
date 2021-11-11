package communication

import (
	"encoding/json"
	"os"
	"path"
)

type UIUpdate struct {
	Receiver string
	Command  string
	Address  string
	Path     string
	Event    string
}

func EmitUpdate(message UIUpdate) error {
	// JSONify message
	messageJson, messageErr := json.Marshal(message)
	if messageErr != nil {
		return messageErr
	}

	// Adding line break to JSON byte output
	// messageJson = append(messageJson, []byte("\n")...)

	// Setting up directory path variables
	homeDir, _ := os.UserHomeDir()
	pluralithDir := path.Join(homeDir, "Pluralith")

	// Create parent directories for path if they don't exist yet
	if mkErr := os.MkdirAll(pluralithDir, 0700); mkErr != nil {
		return mkErr
	}

	// Write to Pluralith UI bus file (WriteFile replaces all file contents)
	if writeErr := os.WriteFile(path.Join(pluralithDir, "pluralith_ui.bus"), messageJson, 0700); writeErr != nil {
		return writeErr
	}

	return nil
}

// // CLI update object structure (CLI receives this)
// {
// 	receiver: "CLI"
// 	command: "confirmed" || "denied"
// 	path: "- working directory -"
// }

// // UI update object structure (CLI sends this)
// {
// 	receiver: "UI"
// 	command: "apply" || "destroy" || "plan"
// 	path: "- working directory -"
// 	event: "complete" || "progress" || "start" || "begin" || "end"
// }
