package ux

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// - - - Custom Spinner Struct - - -

// Custom spinner struct
type spin struct {
	spinMsg    string
	successMsg string
	failMsg    string
	instance   *spinner.Spinner
}

// Defining custom color print functions
var printBlue = color.New(color.FgBlue).SprintFunc()
var printRed = color.New(color.FgRed).SprintFunc()

// Method to instantiate customer spinner
func NewSpinner(spinMsg string, successMsg string, failMsg string) spin {
	// Creating base spinner instance
	instance := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	// Adding spinning message
	instance.Suffix = " " + spinMsg
	// Chaning color
	instance.Color("blue")

	// Returning custom spin instance
	s := spin{spinMsg, successMsg, failMsg, instance}
	return s
}

// Method to start custom spinner
func (s spin) Start() {
	s.instance.Start()
}

// Method to update spinner state on successful completion
func (s spin) Success(customMessage ...string) {
	// Handling custom message if given
	var message string
	if len(customMessage) > 0 {
		message = customMessage[0]
	} else {
		message = s.successMsg
	}
	// Stopping base spinner
	s.instance.Stop()
	// Printing custom success message
	fmt.Printf("%s %s\n", printBlue("✔"), message)
}

// Method to update spinner state on failure
func (s spin) Fail(customMessage ...string) {
	// Handling custom message if given
	var message string
	if len(customMessage) > 0 {
		message = customMessage[0]
	} else {
		message = s.failMsg
	}
	// Stopping base spinner
	s.instance.Stop()
	// Printing custom failure message
	fmt.Printf("%s %s\n", printRed("✖️"), message)
}
