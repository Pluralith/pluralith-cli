package apply

import (
	"fmt"

	"github.com/fatih/color"

	plan "pluralith/pkg/plan"
	stream "pluralith/pkg/stream"
)

// Defining custom color print functions
var printBlue = color.New(color.FgBlue).SprintFunc()

func ApplyMethod(args []string) {
	// Printing custom success message
	fmt.Printf("%s %s\n\n", printBlue("⠿"), "Running apply, confirm it in the Desktop app")
	// Run Pluralith plan routine
	plan.PlanMethod(args)

	stream.StreamOutput(args)

}
