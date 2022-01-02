package auxiliary

import (
	"fmt"
	"os"
	"path/filepath"
)

type Paths struct {
	HomePath    string
	WorkingPath string
	ComDBPath   string
	LockPath    string
}

func (P *Paths) GeneratePaths() error {
	functionName := "GeneratePaths"

	// Get current working directory
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return fmt.Errorf("%v: %w", functionName, workingErr)
	}

	// Get user home directory
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf("%v: %w", functionName, homeErr)
	}

	// Set path parameters
	P.HomePath = homeDir
	P.WorkingPath = workingDir
	P.ComDBPath = filepath.Join(homeDir, "Pluralith", "pluralithComDB.json")
	P.LockPath = filepath.Join(homeDir, "Pluralith", "pluralithLock.json")

	return nil
}

var PathInstance = &Paths{}
