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

func StreamCommand(command string, args []string) error {
	functionName := "StreamCommand"

	// Instantiate spinners
	streamSpinner := ux.NewSpinner("Apply Running", "Apply Completed", "Apply Failed", true)
	// Adapting spinner to destroy command
	if command == "destroy" {
		streamSpinner = ux.NewSpinner("Destroy Running", "Destroy Completed", "Destroy Failed", true)
	}

	// Emit apply begin update to UI
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   command,
		Type:      "begin",
		Path:      auxiliary.StateInstance.WorkingPath,
		Received:  false,
	})

	streamSpinner.Start()

	// Constructing command to execute
	cmd := exec.Command("terraform", args...)

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
			Path:      auxiliary.StateInstance.WorkingPath,
			Received:  false,
		})

		return fmt.Errorf("%v: %w", functionName, cmdErr)
	}

	// Scan for command line updates
	applyScanner := bufio.NewScanner(outStream)
	applyScanner.Split(bufio.ScanLines)

	// While command line scan is running
	for applyScanner.Scan() {
		ProcessTerraformMessage(applyScanner.Text(), command)
	}

	// Pull state with latest attributes
	latestState, pullErr := PullState()
	if pullErr != nil {
		return fmt.Errorf("pulling terraform state failed -> %v: %w", functionName, pullErr)
	}

	// Emit apply end update to UI
	comdb.PushComDBEvent(comdb.ComDBEvent{
		Receiver:  "UI",
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Command:   command,
		Type:      "end",
		Path:      auxiliary.StateInstance.WorkingPath,
		State:     latestState,
		Received:  false,
	})

	streamSpinner.Success()

	ux.PrintFormatted("\n✔ All Done\n", []string{"green", "bold"})

	return nil
}
