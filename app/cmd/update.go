/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	Short: "Check for and install updates for the Pluralith CLI",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		fmt.Print("Current Version: ")
		ux.PrintFormatted(auxiliary.StateInstance.CLIVersion+"\n\n", []string{"bold", "blue"})

		url := "https://api.pluralith.com/v1/dist/download/cli"
		params := map[string]string{"os": runtime.GOOS, "arch": runtime.GOARCH}

		updateUrl, shouldUpdate, checkErr := install.GetGitHubRelease(url, params, auxiliary.StateInstance.CLIVersion)
		if checkErr != nil {
			fmt.Println(fmt.Errorf("failed to get latest release -> %w", checkErr))
		}

		// Get source path of current executable to download update to
		UpdatePath, pathErr := os.Executable()
		if pathErr != nil {
			fmt.Println(fmt.Errorf("failed to get CLI binary source path -> %w", pathErr))
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
			}
		}

		if shouldUpdate {
			if downloadErr := install.DownloadGitHubRelease("Pluralith CLI", updateUrl, UpdatePath); downloadErr != nil {
				fmt.Println(fmt.Errorf("failed to download latest version -> %w", downloadErr))
			}

			// Execute batch script to replace binary after CLI terminates
			if runtime.GOOS == "windows" {
				replaceCmd := exec.Command(filepath.Join(auxiliary.StateInstance.BinPath, "update.bat"))
				if replaceErr := replaceCmd.Start(); replaceErr != nil {
					fmt.Println(fmt.Errorf("failed to run update helper -> %w", replaceErr))
				}
			}
		}
	},
}

// Graph module
var updateGraphModule = &cobra.Command{
	Use:   "graph-module",
	Short: "Strip a given state file of secrets according to config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		components.GraphModule()
	},
}

func init() {
	updateCmd.AddCommand(updateGraphModule)
	rootCmd.AddCommand(updateCmd)
}
