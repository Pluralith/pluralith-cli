package auxiliary

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/ci"
	"runtime"
	"strings"
)

type State struct {
	CLIVersion    string
	HomePath      string
	WorkingPath   string
	PluralithPath string
	BinPath       string
	ComDBPath     string
	LockPath      string
	APIKey        string
	IsWSL         bool
	IsCI          bool
}

func (P *State) GeneratePaths() error {
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
		pathParts := strings.Split(WSLPath, "Pluralith")
		homeDir = pathParts[0]
	}

	// Set path parameters
	P.HomePath = homeDir
	P.WorkingPath = workingDir
	P.PluralithPath = filepath.Join(homeDir, "Pluralith")
	P.BinPath = filepath.Join(P.PluralithPath, "bin")
	P.ComDBPath = filepath.Join(P.PluralithPath, "pluralithComDB.json")
	P.LockPath = filepath.Join(P.PluralithPath, "pluralithLock.json")

	return nil
}

func (P *State) InitPaths() error {
	functionName := "InitPaths"

	// Create parent directories for path if they don't exist yet
	if mkErr := os.MkdirAll(P.BinPath, 0700); mkErr != nil {
		return fmt.Errorf("%v: %w", functionName, mkErr)
	}

	return nil
}

func (P *State) SetAPIKey() error {
	functionName := "SetAPIKey"
	credentialsPath := filepath.Join(P.PluralithPath, "credentials")

	// Check if credentials file exists
	if _, pathErr := os.Stat(credentialsPath); errors.Is(pathErr, os.ErrNotExist) {
		P.APIKey = ""
		return nil
	}

	keyValue, readErr := os.ReadFile(credentialsPath)
	if readErr != nil {
		return fmt.Errorf("%v: %w", functionName, readErr)
	}

	P.APIKey = string(keyValue)
	return nil
}

func (P *State) CheckWSL() string {
	// If OS is some form of Linux
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		// Get kernel version
		versionBytes, versionErr := os.ReadFile("/proc/version")
		if versionErr != nil {
			P.IsWSL = false
			return ""
		}

		versionString := strings.ToLower(string(versionBytes))

		// If version string contains microsoft -> Linux running in WSL
		if strings.Contains(versionString, "microsoft") {
			// Get executable source directory
			ex, err := os.Executable()
			if err != nil {
				fmt.Println("Could not check for WSL")
				P.IsWSL = false
				return ""
			}

			P.IsWSL = true
			return filepath.Dir(ex)
		}
	}

	P.IsWSL = false
	return ""
}

func (P *State) CheckCI() {
	if isCI := ci.CheckEnvVars(); isCI {
		P.IsCI = true
		return
	}

	if isCI := ci.CheckDocker(); isCI {
		P.IsCI = true
		return
	}

	P.IsCI = false
}

var StateInstance = &State{}
