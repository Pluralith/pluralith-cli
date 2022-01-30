package auxiliary

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Paths struct {
	HomePath      string
	WorkingPath   string
	PluralithPath string
	ComDBPath     string
	LockPath      string
	IsWSL         bool
}

func (P *Paths) CheckWSL() string {
	// If OS is some form of Linux
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		// Get executable source directory
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("Could not check for WSL")
			P.IsWSL = false
			return ""
		}

		exPath := filepath.Dir(ex)

		// If source directory contains "/mnt/" -> assume it is WSL
		if strings.Contains(exPath, "/mnt/") {
			P.IsWSL = true
			return exPath
		}
	}

	P.IsWSL = false
	return ""
}

func (P *Paths) GeneratePaths() error {
	functionName := "GeneratePaths"

	// Check for WSL
	WSLPath := P.CheckWSL()

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

	// If it is WSL -> Set homedir to homedir on Windows
	if P.IsWSL {
		pathParts := strings.Split(WSLPath, "/")
		homeDir = strings.Join(pathParts[:5], "/")
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
