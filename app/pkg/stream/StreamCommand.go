package stream

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"pluralith/pkg/auxiliary"
	"pluralith/pkg/comdb"
	"pluralith/pkg/ux"
)

func handleTerraformOutput(jsonString string, command string) error {
	functionName := "handleTerraformOutput"
	// Decode json string to get event type and resource address
	event, decodeErr := DecodeStateStream(jsonString, command)
	if decodeErr != nil {
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
			Command:   event.Command,
			Type:      event.Type,
			Address:   event.Address,
			Message:   event.Message,
			Instances: instances,
			Path:      auxiliary.PathInstance.WorkingPath,
			Received:  false,
		})
	}

	return nil
}

func StreamCommand(command string, args []string) error {
	functionName := "StreamCommand"

	// Instantiate spinners
	streamSpinner := ux.NewSpinner("Apply Running", "Apply Completed", "Apply Failed")
	// Adapting spinner to destroy command
	if command == "destroy" {
		streamSpinner = ux.NewSpinner("Destroy Running", "Destroy Completed", "Destroy Failed")
	}

	// Emit apply begin update to UI
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   command,
		Type:      "begin",
		Instances: make([]interface{}, 0),
		Path:      auxiliary.PathInstance.WorkingPath,
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
		fmt.Println(errorSink.String())

		comdb.PushComDBEvent(comdb.ComDBEvent{
			Receiver:  "UI",
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Command:   command,
			Type:      "failed",
			Error:     errorSink.String(),
			Path:      auxiliary.PathInstance.WorkingPath,
			Received:  false,
		})

		return fmt.Errorf("%v: %w", functionName, cmdErr)
	}

	// Scan for command line updates
	applyScanner := bufio.NewScanner(outStream)
	applyScanner.Split(bufio.ScanLines)

	// While command line scan is running
	for applyScanner.Scan() {
		if scanErr := handleTerraformOutput(applyScanner.Text(), command); scanErr != nil {
			streamSpinner.Fail()
			return fmt.Errorf("scanning terraform json output failed -> %v: %w", functionName, scanErr)
		}
	}

	// Emit apply end update to UI
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   command,
		Type:      "end",
		Instances: make([]interface{}, 0),
		Path:      auxiliary.PathInstance.WorkingPath,
		Received:  false,
	})

	streamSpinner.Success()

	return nil
}
