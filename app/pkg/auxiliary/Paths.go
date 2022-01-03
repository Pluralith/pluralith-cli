package auxiliary

import (
	"fmt"
	"os"
	"path/filepath"
)

type Paths struct {
	HomePath      string
	WorkingPath   string
	PluralithPath string
	ComDBPath     string
	LockPath      string
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
	P.PluralithPath = filepath.Join(homeDir, "Pluralith") 
	P.ComDBPath = filepath.Join(P.PluralithPath, "pluralithComDB.json")
	P.LockPath = filepath.Join(P.PluralithPath, "pluralithLock.json")

	return nil
}

func (P *Paths) InitPaths() error {
	functionName := "InitPaths"

	// Create parent directories for path if they don't exist yet
	if mkErr := os.MkdirAll(P.PluralithPath, 0700); mkErr != nil {
		return fmt.Errorf("%v: %w", functionName, mkErr)
	}

	return nil
}

var PathInstance = &Paths{}
