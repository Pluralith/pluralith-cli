package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/install"
	"pluralith/pkg/install/components"
	"pluralith/pkg/ux"
	"runtime"

	"github.com/spf13/cobra"
)

var winUpdateHelper string = `TIMEOUT 3
DEL pluralith.exe
REN pluralith_update.exe pluralith.exe`

// planOldCmd represents the planOld command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates for the Pluralith CLI and its modules",
	Long:  `Check for and install updates for the Pluralith CLI and its modules`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		fmt.Print("Current Version: ")
		ux.PrintFormatted(auxiliary.StateInstance.CLIVersion+"\n\n", []string{"bold", "blue"})

		url := "https://api.pluralith.com/v1/dist/download/cli"
		params := map[string]string{"os": runtime.GOOS, "arch": runtime.GOARCH}

		updateUrl, shouldUpdate, checkErr := install.GetGitHubRelease(url, params, auxiliary.StateInstance.CLIVersion)
		if checkErr != nil {
			fmt.Println(fmt.Errorf("failed to get latest release -> %w", checkErr))
			return
		}

		// Get source path of current executable to download update to
		UpdatePath, pathErr := os.Executable()
		if pathErr != nil {
			fmt.Println(fmt.Errorf("failed to get CLI binary source path -> %w", pathErr))
			return
		}

		// Add special condition where windows cannot override currently executed binary
		// 1) Write to separate update file
		// 2) Ensure update helper batch script exists
		// 2) Execute batch script with delay that renames and replaces current binary with updated binary
		if runtime.GOOS == "windows" {
			UpdatePath = filepath.Join(auxiliary.StateInstance.BinPath, "pluralith_update.exe")
			UpdateHelperPath := filepath.Join(auxiliary.StateInstance.BinPath, "update.bat")
			helperWriteErr := os.WriteFile(UpdateHelperPath, []byte(winUpdateHelper), 0700)
			if helperWriteErr != nil {
				fmt.Println(fmt.Errorf("failed to create update helper -> %w", helperWriteErr))
				return
			}
		}

		if shouldUpdate {
			if downloadErr := install.DownloadGitHubRelease("Pluralith CLI", updateUrl, UpdatePath); downloadErr != nil {
				fmt.Println(fmt.Errorf("failed to download latest version -> %w", downloadErr))
				return
			}

			// Execute batch script to replace binary after CLI terminates
			if runtime.GOOS == "windows" {
				replaceCmd := exec.Command(filepath.Join(auxiliary.StateInstance.BinPath, "update.bat"))
				if replaceErr := replaceCmd.Start(); replaceErr != nil {
					fmt.Println(fmt.Errorf("failed to run update helper -> %w", replaceErr))
					return
				}
			}
		}
	},
}

// Graph module
var updateGraphModule = &cobra.Command{
	Use:   "graph-module",
	Short: "Install the latest graph module for the Pluralith CLI",
	Long:  `Install the latest graph module for the Pluralith CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		components.GraphModule()
	},
}

func init() {
	updateCmd.AddCommand(updateGraphModule)
	rootCmd.AddCommand(updateCmd)
}
