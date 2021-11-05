package auxiliary

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"pluralith/helpers"
)

func GenerateJson(planPath string) (string, error) {
	// Manually parsing arg (due to cobra lacking a feature)
	// parsedArgs, parsedArgMap := helpers.ParseArgs(args, pluralithPlanArgs)
	// Getting value of -out flag
	// planOut := parsedArgMap["out"]

	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{"show", "-json", planPath})...)

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

	strippedFile, stripErr := helpers.StripSecrets(outputSink.String(), []string{}, "gatewatch")
	if stripErr == nil {
		ioutil.WriteFile("pluralith.state.stripped", []byte(strippedFile), 0644)
	}

	// Returning location of execution plan
	return strippedFile, nil
}
