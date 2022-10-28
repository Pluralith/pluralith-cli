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

func PadEventLogs(address string, newEvent []string, eventLog *[][]string, currentPadding *int) {
	// Adapt padding if necessary
	if len(address) > *currentPadding {
		*currentPadding = len(address) + 1
	}

	// Append new event
	*eventLog = append(*eventLog, newEvent)

	// Calculate padding for individual log lines
	for _, log := range *eventLog {
		// If current log's resource address is shorter than current padding -> Increase padding
		if len(log[1]) < *currentPadding {
			paddingLength := *currentPadding - len(log[1]) // Calculate new padding for current log line
			log[2] = ""                                    // Reset previous padding

			// Fill padding string ([2] in log slice) with padding spaces
			for padding := 0; padding <= paddingLength; padding++ {
				log[2] += " "
			}
		}
	}
}

func StreamCommand(command string, tfArgs []string) error {
	functionName := "StreamCommand"

	// Set command mode for prints
	commandMode := "Apply"
	if command == "destroy" {
		commandMode = "Destroy"
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

	// Constructing command to execute
	cmd := exec.Command("terraform", tfArgs...)

	// Define sinks for std data
	var errorSink bytes.Buffer

	// Redirect command std data
	cmd.Stderr = &errorSink

	// Initiate standard output pipe
	outStream, outErr := cmd.StdoutPipe()
	if outErr != nil {
		ux.PrintFormatted("\n  ✘ ", []string{"bold", "red"})
		fmt.Println(commandMode + " Failed")

		return fmt.Errorf("%v: %w", functionName, outErr)
	}

	// Run terraform command
	cmdErr := cmd.Start()
	if cmdErr != nil {
		ux.PrintFormatted("\n  ✘ ", []string{"bold", "red"})
		fmt.Println(commandMode + " Failed")

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

	// Deactivate cursor
	fmt.Print("\033[?25l\n")

	// While command line scan is running
	for applyScanner.Scan() {
		event := ProcessTerraformMessage(applyScanner.Text(), command)

		// If address is given -> Resource event
		if event.Address != "" {
			// Emit current event update to UI
			comdb.PushComDBEvent(comdb.ComDBEvent{
				Receiver:  "UI",
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Command:   event.Command,
				Type:      event.ParsedType,
				Address:   event.Address,
				Message:   event.Message,
				Path:      auxiliary.StateInstance.WorkingPath,
				Received:  false,
			})
		}

		if event.Type != "version" && event.Type != "planned_change" && event.Type != "outputs" {
			fmt.Println("    " + event.Message)
		}

		if event.Type == "outputs" {
			ux.PrintFormatted("\n  → ", []string{"bold", "blue"})
			ux.PrintFormatted(event.Message, []string{"bold"})
			for outputName, outputObject := range event.Outputs {
				ux.PrintFormatted("\n    → ", []string{"bold"})
				fmt.Println(outputName+" =", outputObject.(map[string]interface{})["value"])
			}
		}
	}

	// Pull state with latest attributes
	latestState, pullErr := PullState()
	if pullErr != nil {
		ux.PrintFormatted("\n  ✘ ", []string{"bold", "red"})
		fmt.Println(commandMode + " Failed")

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

	ux.PrintFormatted("\n  ✔ ", []string{"bold", "blue"})
	fmt.Println(commandMode + " Completed")

	ux.PrintFormatted("\n✔ All Done\n", []string{"green", "bold"})

	// Activate cursor
	fmt.Print("\033[?25h")

	return nil
}
