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

func ExecuteTerraform(command string, args []string, spinner bool, stdOut bool, stdErr bool) (string, int) {
	// Capitalizing command string for UX
	titledCommand := strings.Title(command)
	// Instantiating new spinner
	terraformSpinner := ux.NewSpinner("Running Terraform "+titledCommand, titledCommand+" done", titledCommand+" failed")
	if spinner {
		terraformSpinner.Start()
	}
	
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
		// Getting error code
		if exit, ok := err.(*exec.ExitError); ok {
			// Handling UX according to passed args
			if spinner {
				terraformSpinner.Fail()
			}
			if stdErr {
				fmt.Println(errorSink.String())
			}

			// Returning stdout as string aswell as whatever error code the command yielded
			return outputSink.String(), exit.ExitCode()
		}
	}

	if spinner {
		terraformSpinner.Success()
	}

	// Returning stdout as string aswell as error dode 0
	return outputSink.String(), 0
}
