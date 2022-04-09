package terraform

import (
	"fmt"
	"os"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/comdb"
	"pluralith/pkg/stream"
	"pluralith/pkg/ux"
	"strings"
	"time"
)

func RunApply(command string, args []string) error {
	functionName := "RunApply"

	// Define terraform args used by this command
	pluralithArgs := make(map[string]string)
	pluralithArgs["auto-approve"] = "true"
	pluralithArgs["json"] = "true"
	pluralithArgs["input"] = "false"
	if command == "destroy" {
		pluralithArgs["destroy"] = ""
	}

	// Manually parse arg (due to cobra lacking a feature)
	parsedArgs := auxiliary.ParseArgs(args, pluralithArgs)

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

	var confirm bool
	var watchErr error

	// Watch for updates from UI and wait for confirmation
	if auxiliary.StateInstance.IsWSL { // Watch ComDB with loop due to missing inotify support in WSL 2
		confirm, watchErr = comdb.WatchComDBFallback()
		if watchErr != nil {
			return fmt.Errorf("instantiating ComDB watcher failed -> %v: %w", functionName, watchErr)
		}
	} else { // Use actual file watcher for everything else
		confirm, watchErr = comdb.WatchComDB()
		if watchErr != nil {
			return fmt.Errorf("instantiating ComDB watcher failed -> %v: %w", functionName, watchErr)
		}
	}

	// Stream apply command output
	if confirm {
		confirmSpinner.Success()

		// Adapt command string if plan
		if command == "plan" {
			command = "apply"
		}

		streamErr := stream.StreamCommand(command, parsedArgs)
		if streamErr != nil {
			return fmt.Errorf("streaming terraform command output failed -> %v: %w", functionName, streamErr)
		}
	} else {
		confirmSpinner.Fail()
	}

	return nil
}
