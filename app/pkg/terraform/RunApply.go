package terraform

import (
	"os"
	"pluralith/pkg/comdb"
	"pluralith/pkg/stream"
	"pluralith/pkg/ux"
	"time"
)

func RunApply(command string, args []string) error {
	// Get working directory
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return workingErr
	}

	// Instantiate spinner
	confirmSpinner := ux.NewSpinner(
		RunMessages[command].([]string)[1],
		RunMessages[command].([]string)[2],
		RunMessages[command].([]string)[3],
	)

	// Emit confirm event
	comdb.PushComDBEvent(comdb.Event{
		Receiver:   "UI",
		Timestamp:  time.Now().Unix(),
		Command:    command, // UI can only mark as received when command is "apply" for some reason
		Type:       "confirm",
		Address:    "",
		Attributes: make(map[string]interface{}),
		Path:       workingDir,
		Received:   false,
	})

	confirmSpinner.Start()

	// Watch for updates from UI and wait for confirmation
	confirm, watchErr := comdb.WatchComDB()
	if watchErr != nil {
		return watchErr
	}

	// Stream apply command output
	if confirm {
		confirmSpinner.Success()
		streamErr := stream.StreamCommand(command, args)
		if streamErr != nil {
			return streamErr
		}
	} else {
		confirmSpinner.Fail()
	}

	return nil
}
