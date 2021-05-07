package helpers

import (
	"os/exec"
	"runtime"
	"time"

	"pluralith/ux"
)

// - - - Code to launch Pluralith UI cross-platform - - -

// Function to run OS specific launch command
func runOsCommand(command []string) {
	// Instantiating new custom spinner
	spinner := ux.NewSpinner("Launching Pluralith...", "Pluralith Running\n", "Failed to launch Pluralith\n")
	spinner.Start()

	// Creating command to launch Pluralith on given OS
	cmd := exec.Command(command[0], command[1:]...)
	// Handling success and failure cases for terminal command
	// Adding slight delay to debounce for UI to get there
	if err := cmd.Run(); err != nil {
		time.Sleep(300 * time.Millisecond)
		spinner.Fail()
	} else {
		time.Sleep(300 * time.Millisecond)
		spinner.Success()
	}
}

// Function to launch Pluralith UI
func LaunchPluralith() {
	// Running terminal command to launch application based on current OS
	switch os := runtime.GOOS; os {
	case "windows":
		runOsCommand([]string{"command", "and", "arguments"})
	case "darwin":
		runOsCommand([]string{"open", "-a", "Pluralith"})
	default:
		runOsCommand([]string{"command", "and", "arguments"})
	}
}
