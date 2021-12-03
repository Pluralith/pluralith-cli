package helpers

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

// - - - Code to execute terraform commands - - -

func ExecuteTerraform(command string, args []string, stdOut bool) (string, error) {
	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{command}, args...)...)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink
	cmd.Stdin = os.Stdin

	// Printing std output to console if arg set to true
	if stdOut {
		cmd.Stdout = os.Stdout
	}

	// Running terraform command
	if err := cmd.Run(); err != nil {
		// Handling error
		return errorSink.String(), errors.New("terraform command failed")
	}

	// Returning stdout as string
	return outputSink.String(), nil
}
