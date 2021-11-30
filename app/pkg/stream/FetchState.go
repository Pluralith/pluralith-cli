package stream

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"pluralith/pkg/auxiliary"
)

func FetchState(address string) (map[string]interface{}, error) {
	// Constructing command to execute
	cmd := exec.Command("terraform", "state", "pull")

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer
	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink

	// Running terraform command
	if err := cmd.Run(); err != nil {
		fmt.Println(errorSink.String())
		// Handling error
		return make(map[string]interface{}), errors.New("terraform command failed")
	}

	// Parse state received from terraform state pull
	parsedState, parseErr := auxiliary.ParseJson(outputSink.String()) // FAILS ON DESTROY -> FIX
	if parseErr != nil {
		return make(map[string]interface{}), parseErr
	}

	return parsedState, nil
}
