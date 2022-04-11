package stream

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"pluralith/pkg/auxiliary"
	"pluralith/pkg/comdb"
	"pluralith/pkg/ux"

	"github.com/fatih/color"
)

func StreamCommand(command string, tfArgs []string) error {
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

	// streamSpinner.Start()

	// Constructing command to execute
	cmd := exec.Command("terraform", tfArgs...)

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
		// streamSpinner.Fail()
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

	// progressWriter := color.Output
	eventLog := ""
	eventCount := 0

	errorCount := 0
	errorPrint := color.New(color.Bold, color.FgHiRed)

	commandCount := 0
	commandMode := "Completed"
	commandModePrint := color.New(color.Bold, color.FgHiGreen)

	if command == "destroy" {
		commandMode = "Destroyed"
		commandModePrint = color.New(color.Bold, color.FgHiBlue)
	}

	// blueStyle := color.New(color.Bold, color.FgHiBlue)
	// redStyle := color.New(color.Bold, color.FgHiRed)
	// greenStyle :=

	commandModePrint.Sprint(strconv.Itoa(commandCount))

	// Deactivate cursor
	fmt.Print("\033[?25l")

	ux.PrintFormatted("  → ", []string{"bold", "blue"})
	// commandCountString :=
	fmt.Printf("Running → %s %s / %s Errored", commandModePrint.Sprint(strconv.Itoa(commandCount)), commandMode, errorPrint.Sprint(strconv.Itoa(errorCount)))

	// While command line scan is running
	for applyScanner.Scan() {
		event := ProcessTerraformMessage(applyScanner.Text(), command)
		var eventString string

		if event.Type == "complete" {
			commandCount += 1
			eventString = fmt.Sprintf("%s %s %s", commandModePrint.Sprint("    ✔ "), event.Address, commandModePrint.Sprint(commandMode))
		}

		if event.Type == "errored" {
			errorCount += 1
			eventString = fmt.Sprintf("%s %s %s", errorPrint.Sprint("    ✘ "), event.Address, errorPrint.Sprint("Errored"))
		}

		// fmt.Println(messageString)
		if event.Address != "" && eventString != "" {
			for line := 0; line <= eventCount; line++ {
				fmt.Printf("\033[A")
			}

			fmt.Println()
			eventCount += 1

			ux.PrintFormatted("  → ", []string{"bold", "blue"})
			fmt.Printf("Running → %s %s / %s Errored", commandModePrint.Sprint(strconv.Itoa(commandCount)), commandMode, errorPrint.Sprint(strconv.Itoa(errorCount)))
			fmt.Printf("\033F")

			eventLog += "\n" + eventString
			fmt.Print(eventLog)
		}
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

	// streamSpinner.Success()

	ux.PrintFormatted("\n\n✔ All Done\n", []string{"green", "bold"})

	// Activate cursor
	fmt.Print("\033[?25h")

	return nil
}
