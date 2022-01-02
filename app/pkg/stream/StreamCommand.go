package stream

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"pluralith/pkg/comdb"
	"pluralith/pkg/ux"
)

func StreamCommand(command string, args []string) error {
	functionName := "StreamCommand"

	// Instantiate spinners
	streamSpinner := ux.NewSpinner("Apply Running", "Apply Completed", "Apply Failed")
	// Adapting spinner to destroy command
	if command == "destroy" {
		streamSpinner = ux.NewSpinner("Destroy Running", "Destroy Completed", "Destroy Failed")
	}

	// Get working directory for update emission
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return fmt.Errorf("%v: %w", functionName, workingErr)
	}

	// Emit apply begin update to UI
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   command,
		Type:      "begin",
		Instances: make([]interface{}, 0),
		Path:      workingDir,
		Received:  false,
	})

	streamSpinner.Start()

	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{"apply", "-lock=false"}, args...)...)

	// Define sinks for std data
	var errorSink bytes.Buffer

	// Redirect command std data
	cmd.Stderr = &errorSink

	// Initiate standard output pipe
	outStream, outErr := cmd.StdoutPipe()
	if outErr != nil {
		streamSpinner.Fail()
		return fmt.Errorf("%v: %w", functionName, outErr)
	}

	// Run terraform command
	cmdErr := cmd.Start()
	if cmdErr != nil {
		streamSpinner.Fail()
		return fmt.Errorf("%v: %w", functionName, cmdErr)
	}

	// Scan for command line updates
	applyScanner := bufio.NewScanner(outStream)
	applyScanner.Split(bufio.ScanLines)

	// While command line scan is running
	for applyScanner.Scan() {
		// Get current line json string
		jsonString := applyScanner.Text()
		// Decode json string to get event type and resource address
		event, decodeErr := DecodeStateStream(jsonString)
		if decodeErr != nil {
			streamSpinner.Fail()
			return fmt.Errorf("%v: %w", functionName, decodeErr)
		}

		// If address is given -> Resource event
		if event.Address != "" {
			var instances []interface{}

			// If event complete -> Fetch resource instances with attributes
			if event.Type == "complete" {
				fetchedState, fetchErr := PullState(event.Address)
				if fetchErr != nil {
					return fmt.Errorf("pulling terraform state failed -> %v: %w", functionName, fetchErr)
				}

				instances = FetchResourceInstances(event.Address, fetchedState)
			}

			// // Emit current event update to UI
			comdb.PushComDBEvent(comdb.ComDBEvent{
				Receiver:  "UI",
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Command:   command,
				Type:      event.Type,
				Address:   event.Address,
				Message:   event.Message,
				Instances: instances,
				Path:      workingDir,
				Received:  false,
			})
		}
	}

	// Emit apply start update to UI
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   command,
		Type:      "end",
		Instances: make([]interface{}, 0),
		Path:      workingDir,
		Received:  false,
	})

	streamSpinner.Success()

	return nil
}
