package auxiliary

import (
	"os/exec"
	"runtime"
	"time"
	"path/filepath"
	"os"
	"fmt"

	"pluralith/pkg/ux"
)

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
func LaunchPluralith() error {
	functionName := "LaunchPluralith"

	// Get homedir
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf("%v: %w", functionName, homeErr)
	}
	// Running terminal command to launch application based on current OS
	switch os := runtime.GOOS; os {
	case "windows":
		runOsCommand([]string{filepath.Join(homeDir, "AppData", "Local", "Programs", "pluralith", "Pluralith.exe")})
	case "darwin":
		runOsCommand([]string{"open", "-a", "Pluralith"})
	default:
		runOsCommand([]string{"command", "and", "arguments"})
	}

	return nil
}
