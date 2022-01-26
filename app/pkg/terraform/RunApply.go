package terraform

import (
	"fmt"
	"os"
	"pluralith/pkg/comdb"
	"pluralith/pkg/stream"
	"pluralith/pkg/ux"
	"strings"
	"time"
)

func RunApply(command string, args []string) error {
	functionName := "RunApply"

	ux.PrintFormatted("→", []string{"blue", "bold"})
	ux.PrintFormatted(strings.Join([]string{" ", strings.Title(command)}, ""), []string{"white", "bold"})
	fmt.Println("")

	// Get working directory
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return fmt.Errorf("%v: %w", functionName, workingErr)
	}

	// Instantiate spinner
	confirmSpinner := ux.NewSpinner(RunMessages[command].([]string)[1], RunMessages[command].([]string)[2], RunMessages[command].([]string)[3], true)

	// Emit confirm event
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   command, // UI can only mark as received when command is "apply" for some reason
		Type:      "confirm",
		Path:      workingDir,
		Received:  false,
	})

	confirmSpinner.Start()

	// Watch for updates from UI and wait for confirmation
	confirm, watchErr := comdb.WatchComDB()
	if watchErr != nil {
		return fmt.Errorf("instantiating ComDB watcher failed -> %v: %w", functionName, watchErr)
	}

	// Stream apply command output
	if confirm {
		confirmSpinner.Success()

		// Adapt command string if plan
		if command == "plan" {
			command = "apply"
		}

		streamErr := stream.StreamCommand(command, args)
		if streamErr != nil {
			return fmt.Errorf("streaming terraform command output failed -> %v: %w", functionName, streamErr)
		}
	} else {
		confirmSpinner.Fail()
	}

	return nil
}
