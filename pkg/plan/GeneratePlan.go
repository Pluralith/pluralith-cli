package plan

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	auxiliary "pluralith/pkg/auxiliary"
)

// Defining command args/flags
var pluralithPlanArgs = []string{"-show-output", "-s"}

func GeneratePlan(args []string) (string, error) {
	// Manually parsing arg (due to cobra lacking a feature)
	parsedArgs, parsedArgMap := auxiliary.ParseArgs(args, pluralithPlanArgs)
	// Getting value of -out flag
	planOut := parsedArgMap["out"]

	// If no value is given for -out, replace it with standard ./pluralith
	if planOut == "" {
		planOut = "./pluralith.plan"
		parsedArgs = append(parsedArgs, "-out", planOut)
	}
	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{"plan"}, parsedArgs...)...)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return errorSink.String(), errors.New("terraform command failed")
	}

	// Returning location of execution plan
	return planOut, nil
}
