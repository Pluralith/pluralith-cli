package apply

import (
	"fmt"

	"pluralith/pkg/auxiliary"
	"pluralith/pkg/communication"
	plan "pluralith/pkg/plan"
	stream "pluralith/pkg/stream"
	ux "pluralith/pkg/ux"
)

func ApplyMethod(args []string) error {
	// Instantiate wait spinner
	confirmationSpinner := ux.NewSpinner("Waiting for Confirmation ⇢ Confirm apply in the Pluralith UI", "Apply Confirmed", "Apply Canceled")
	// Print running message
	ux.PrintFormatted("⠿", []string{"blue"})
	fmt.Println(" Running apply ⇢ Confirm it in the Pluralith UI\n")

	// Manually parse arg (due to cobra lacking a feature)
	parsedArgs, parsedArgMap := auxiliary.ParseArgs(args, []string{})

	// Add necessary flags if not already given
	if parsedArgMap["auto-approve"] == "" {
		parsedArgs = append(parsedArgs, "-auto-approve")
	}
	if parsedArgMap["json"] == "" {
		parsedArgs = append(parsedArgs, "-json")
	}

	// Run Pluralith plan routine
	planPath, planErr := plan.PlanMethod([]string{""}, true)
	if planErr != nil {
		return planErr
	}

	// Add plan path to arguments to run apply on already created execution plan
	parsedArgs = append(parsedArgs, planPath)

	fmt.Println()
	confirmationSpinner.Start()
	// Watch for updates from UI and wait for confirmation
	confirm, watchErr := communication.WatchForUpdates()
	if watchErr != nil {
		return watchErr
	}

	// Stream apply command output
	if confirm {
		confirmationSpinner.Success()
		streamErr := stream.StreamCommand(parsedArgs, false)
		if streamErr != nil {
			return streamErr
		}
	} else {
		confirmationSpinner.Fail()
	}

	return nil
}
