package helpers

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Function to launch Pluralith UI
func LaunchPluralith() {
	// Running terminal command to launch application based on current OS
	switch os := runtime.GOOS; os {
	case "windows":
		fmt.Println("Pluralith on Windows")
	case "darwin":
		cmd := exec.Command("open", "-a", "Pluralith")
		cmd.Run()
	default:
		fmt.Println("Pluralith on Linux")
	}
}
