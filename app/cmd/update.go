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

// planOldCmd represents the planOld command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates for the Pluralith CLI and its modules",
	Long:  `Check for and install updates for the Pluralith CLI and its modules`,
	Run: func(cmd *cobra.Command, args []string) {
		var unixTargetPath string = "/usr/local/bin"
		if auxiliary.StateInstance.IsWSL {
			unixTargetPath = auxiliary.StateInstance.BinPath
		}

		var unixUpdateHelper string = fmt.Sprintf("sleep 2\nrm %s/pluralith\nmv %s/pluralith_update %s/pluralith", unixTargetPath, auxiliary.StateInstance.BinPath, unixTargetPath)
		var winUpdateHelper string = fmt.Sprintf("TIMEOUT /nobreak /t 2 \nDEL %s\\pluralith.exe \nREN %s\\pluralith_update.exe pluralith.exe", auxiliary.StateInstance.BinPath, auxiliary.StateInstance.BinPath)

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

		// Update process
		var UpdatePath string
		var UpdateHelperPath string
		// 1) Write to separate update file
		// 2) Ensure update helper script exists
		// 2) Execute batch script with delay that renames and replaces current binary with updated binary
		if runtime.GOOS == "windows" {
			UpdatePath = filepath.Join(auxiliary.StateInstance.BinPath, "pluralith_update.exe")
			UpdateHelperPath = filepath.Join(auxiliary.StateInstance.BinPath, "update.bat")
			helperWriteErr := os.WriteFile(UpdateHelperPath, []byte(winUpdateHelper), 0700)
			if helperWriteErr != nil {
				fmt.Println(fmt.Errorf("failed to create update helper -> %w", helperWriteErr))
				return
			}
		} else {
			UpdatePath = filepath.Join(auxiliary.StateInstance.BinPath, "pluralith_update")
			UpdateHelperPath = filepath.Join(auxiliary.StateInstance.BinPath, "update.sh")
			helperWriteErr := os.WriteFile(UpdateHelperPath, []byte(unixUpdateHelper), 0700)
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

			// Depending on os -> run script at path as bash
			var replaceCmd *exec.Cmd
			if runtime.GOOS == "windows" {
				replaceCmd = exec.Command("CMD", "/C", UpdateHelperPath)
			} else {
				replaceCmd = exec.Command("bash", UpdateHelperPath)
			}

			// Execute helper script to replace binary after CLI terminates
			if replaceErr := replaceCmd.Start(); replaceErr != nil {
				fmt.Println(fmt.Errorf("failed to run update helper -> %w", replaceErr))
				return
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
		ux.PrintHead()

		fmt.Print("Installing Latest ")
		ux.PrintFormatted("Graph Module\n\n", []string{"bold", "blue"})

		components.GraphModule()
	},
}

func init() {
	updateCmd.AddCommand(updateGraphModule)
	rootCmd.AddCommand(updateCmd)
}
