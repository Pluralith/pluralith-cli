package auxiliary

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/ci"

	"runtime"
	"strings"
)

type State struct {
	PluralithConfig PluralithConfig
	CLIVersion      string
	HomePath        string
	WorkingPath     string
	PluralithPath   string
	BinPath         string
	ComDBPath       string
	LockPath        string
	APIKey          string
	IsWSL           bool
	IsCI            bool
	Branch          string
	TerraformInit   bool
	Infracost       bool
}

// Produce relevant paths to be used across the application
func (S *State) GeneratePaths() error {
	functionName := "GeneratePaths"

	// Check for WSL
	WSLPath := S.CheckWSL()

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
	if S.IsWSL {
		pathParts := strings.Split(WSLPath, "Pluralith")
		homeDir = pathParts[0]
	}

	// Set path parameters
	S.HomePath = homeDir
	S.WorkingPath = workingDir
	S.PluralithPath = filepath.Join(homeDir, "Pluralith")
	S.BinPath = filepath.Join(S.PluralithPath, "bin")
	S.ComDBPath = filepath.Join(S.PluralithPath, "pluralithComDB.json")
	S.LockPath = filepath.Join(S.PluralithPath, "pluralithLock.json")

	return nil
}

// Make sure important directory structures exist
func (S *State) InitPaths() error {
	functionName := "InitPaths"

	// Create parent directories for path if they don't exist yet
	if mkErr := os.MkdirAll(S.BinPath, 0700); mkErr != nil {
		return fmt.Errorf("%v: %w", functionName, mkErr)
	}

	return nil
}

// Set API key in credentials file
func (S *State) SetAPIKey() error {
	functionName := "SetAPIKey"
	credentialsPath := filepath.Join(S.PluralithPath, "credentials")

	// Check if credentials file exists
	if _, pathErr := os.Stat(credentialsPath); errors.Is(pathErr, os.ErrNotExist) {
		S.APIKey = ""
		return nil
	}

	keyValue, readErr := os.ReadFile(credentialsPath)
	if readErr != nil {
		return fmt.Errorf("%v: %w", functionName, readErr)
	}

	S.APIKey = string(keyValue)
	return nil
}

// Check if running on WSL if GOOS is linux
func (S *State) CheckWSL() string {
	// If OS is some form of Linux
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		// Get kernel version
		versionBytes, versionErr := os.ReadFile("/proc/version")
		if versionErr != nil {
			S.IsWSL = false
			return ""
		}

		versionString := strings.ToLower(string(versionBytes))

		// If version string contains microsoft -> Linux running in WSL
		if strings.Contains(versionString, "microsoft") {
			// Get executable source directory
			ex, err := os.Executable()
			if err != nil {
				fmt.Println("Could not check for WSL")
				S.IsWSL = false
				return ""
			}

			S.IsWSL = true
			return filepath.Dir(ex)
		}
	}

	S.IsWSL = false
	return ""
}

// Check if running in any of the known CI environments
func (S *State) CheckCI() {
	if isCI := ci.CheckEnvVars(); isCI {
		S.IsCI = true
		return
	}

	if isCI := ci.CheckDocker(); isCI {
		S.IsCI = true
		return
	}

	S.IsCI = false
}

// Detect branch name
func (S *State) GetBranch() {
	S.Branch = ci.GetBranch()
}

// Check if Terraform has been initialized
func (S *State) CheckTerraformInit() {
	if _, statErr := os.Stat(filepath.Join(S.WorkingPath, ".terraform")); !os.IsNotExist(statErr) {
		S.TerraformInit = true
	} else {
		S.TerraformInit = false
	}
}

// Check if Infracost is installed
func (S *State) CheckInfracost() {
	verifyCmd := exec.Command("infracost", "--version")

	if verifyErr := verifyCmd.Run(); verifyErr != nil {
		S.Infracost = false
		return
	}

	S.Infracost = true
}

var StateInstance = &State{}
