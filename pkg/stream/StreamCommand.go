package stream

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"pluralith/pkg/communication"
	"pluralith/pkg/ux"
)

func StreamCommand(args []string, isDestroy bool) error {
	// Instantiate spinners
	streamSpinner := ux.NewSpinner("Apply Running", "Apply Completed", "Apply Failed")
	command := "apply"
	// Adapting spinner to destroy command
	if isDestroy {
		streamSpinner = ux.NewSpinner("Destroy Running", "Destroy Completed", "Destroy Failed")
		command = "destroy"
	}

	// Get working directory for update emission
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return workingErr
	}

	// Emit apply begin update to UI
	communication.EmitUpdate(communication.UIUpdate{
		Receiver: "UI",
		Command:  command,
		Address:  "",
		Path:     workingDir,
		Event:    "begin",
	})

	streamSpinner.Start()
	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{"apply"}, args...)...)

	// Define sinks for std data
	var errorSink bytes.Buffer

	// Redirect command std data
	cmd.Stderr = &errorSink

	// Initiate standard output pipe
	outStream, outErr := cmd.StdoutPipe()
	if outErr != nil {
		streamSpinner.Fail()
		return outErr
	}

	// Run terraform command
	cmdErr := cmd.Start()
	if cmdErr != nil {
		streamSpinner.Fail()
		return cmdErr
	}

	// Scan for command line updates
	applyScanner := bufio.NewScanner(outStream)
	applyScanner.Split(bufio.ScanLines)

	// While command line scan is running
	for applyScanner.Scan() {
		// Get current line json string
		jsonString := applyScanner.Text()
		// Decode json string to get event type and resource address
		event, address, decodeErr := DecodeStateStream(jsonString)
		if decodeErr != nil {
			streamSpinner.Fail()
			return decodeErr
		}

		// If address is given
		if address != "" {
			// Fetch current tfstate from state file and strip secrets
			fetchedState, fetchErr := FetchState(address, isDestroy)
			if fetchErr != nil {
				return fetchErr
			}
			FetchResourceAttributes(fetchedState)

			// NOT NECESSARY -> Update plan json and UI will watch those file changes
			// // Emit current event update to UI
			communication.EmitUpdate(communication.UIUpdate{
				Receiver: "UI",
				Command:  command,
				Address:  address,
				Path:     workingDir,
				Event:    strings.Split(event, "_")[1],
			})
		}
	}

	// Emit apply start update to UI
	communication.EmitUpdate(communication.UIUpdate{
		Receiver: "UI",
		Command:  command,
		Address:  "",
		Path:     workingDir,
		Event:    "end",
	})

	streamSpinner.Success()

	return nil
}
