package stream

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	communication "pluralith/pkg/communication"
	ux "pluralith/pkg/ux"
)

func StreamCommand(args []string, isDestroy bool) error {
	// Instantiate spinners
	streamSpinner := ux.NewSpinner("Apply Running", "Apply Completed", "Apply Failed")
	// Adapting spinner to destroy command
	if isDestroy {
		streamSpinner = ux.NewSpinner("Destroy Running", "Destroy Completed", "Destroy Failed")
	}

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

	// Get working directory for update emission
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return workingErr
	}

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

			// Determing command type for update message
			commandType := "apply"
			if isDestroy {
				commandType = "destroy"
			}

			// Construct update message object
			update := communication.UIUpdate{
				Receiver: "UI",
				Command:  commandType,
				Address:  address,
				Path:     workingDir,
				Event:    strings.Split(event, "_")[1],
			}

			// Call function to emit update
			communication.EmitUpdate(update)
		}
	}

	streamSpinner.Success()

	return nil
}
