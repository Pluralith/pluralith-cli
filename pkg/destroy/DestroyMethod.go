package destroy

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/communication"
	"pluralith/pkg/plan"
	"pluralith/pkg/stream"
	"pluralith/pkg/ux"
)

func DestroyMethod(args []string) error {
	// Instantiate wait spinner
	confirmationSpinner := ux.NewSpinner("Waiting for Confirmation ⇢ Confirm destroy in the Pluralith UI", "Destroy Confirmed", "Destroy Canceled")
	// Print running message
	ux.PrintFormatted("⠿", []string{"blue"})
	fmt.Println(" Running destroy ⇢ Confirm it in the Pluralith UI\n")

	// Manually parsing arg (due to cobra lacking a feature)
	parsedArgs, parsedArgMap := auxiliary.ParseArgs(args, []string{})

	// Add necessary flags if not already given
	if parsedArgMap["auto-approve"] == "" {
		parsedArgs = append(parsedArgs, "-auto-approve")
	}
	if parsedArgMap["json"] == "" {
		parsedArgs = append(parsedArgs, "-json")
	}

	// Run Pluralith plan routine
	planPath, planErr := plan.PlanMethod([]string{"-destroy"}, true)
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

	// Stream apply command output (with destroy execution plan)
	if confirm {
		confirmationSpinner.Success()
		streamErr := stream.StreamCommand(parsedArgs, true)
		if streamErr != nil {
			return streamErr
		}
	} else {
		confirmationSpinner.Fail()
	}

	return nil
}
