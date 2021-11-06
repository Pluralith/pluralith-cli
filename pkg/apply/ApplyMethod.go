package apply

import (
	"fmt"

	"pluralith/pkg/auxiliary"
	plan "pluralith/pkg/plan"
	stream "pluralith/pkg/stream"
	ux "pluralith/pkg/ux"
)

func ApplyMethod(args []string) {
	// Print custom success message
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
		fmt.Println(planErr)
	}

	// Add plan path to arguments to run apply on already created execution plan
	parsedArgs = append(parsedArgs, planPath)
	// Stream apply command output
	stream.StreamCommand(parsedArgs, false)
}
