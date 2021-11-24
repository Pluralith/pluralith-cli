package terraform

import (
	"pluralith/pkg/comdb"
	"pluralith/pkg/stream"
	"pluralith/pkg/ux"
)

func RunApply(command string, args []string) error {
	// Instantiate spinner
	confirmSpinner := ux.NewSpinner(
		RunMessages[command].([]string)[1],
		RunMessages[command].([]string)[2],
		RunMessages[command].([]string)[3],
	)

	confirmSpinner.Start()

	// Watch for updates from UI and wait for confirmation
	confirm, watchErr := comdb.WatchComDB()
	if watchErr != nil {
		return watchErr
	}

	// Stream apply command output
	if confirm {
		confirmSpinner.Success()
		streamErr := stream.StreamCommand(args, false)
		if streamErr != nil {
			return streamErr
		}
	} else {
		confirmSpinner.Fail()
	}

	return nil
}
