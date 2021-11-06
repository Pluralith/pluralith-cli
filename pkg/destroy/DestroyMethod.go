package destroy

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/plan"
	"pluralith/pkg/stream"
	"pluralith/pkg/ux"
)

func DestroyMethod(args []string) {
	// Printing custom success message
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
		fmt.Println(planErr)
	}

	// Add plan path to arguments to run apply on already created execution plan
	parsedArgs = append(parsedArgs, planPath)
	// Stream apply command output (with destroy execution plan)
	stream.StreamCommand(parsedArgs, true)
}
