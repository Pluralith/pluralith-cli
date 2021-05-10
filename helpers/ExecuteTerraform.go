package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"pluralith/ux"
	"strings"
)

// - - - Code to execute terraform commands - - -

func ExecuteTerraform(command string, args []string, stdOut bool) (string, int) {
	// Capitalizing command string for UX
	titledCommand := strings.Title(command)
	// Instantiating new spinner
	terraformSpinner := ux.NewSpinner("Running "+titledCommand, titledCommand+" done", titledCommand+" failed")
	terraformSpinner.Start()
	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{command}, args...)...)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink
	cmd.Stdin = os.Stdin

	// Running terraform command
	if err := cmd.Run(); err != nil {
		// Getting error code
		if exit, ok := err.(*exec.ExitError); ok {
			// Stopping spinner and printing stderror to console
			terraformSpinner.Fail()
			fmt.Println(errorSink.String())
			// Returning stdout as string aswell as whatever error code the command yielded
			return outputSink.String(), exit.ExitCode()
		}
	}

	terraformSpinner.Success()

	// If command is not set to silent -> Print stdout
	if stdOut {
		fmt.Println(outputSink.String())
	}

	// Returning stdout as string aswell as error dode 0
	return outputSink.String(), 0
}
