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

		checkSpinner := ux.NewSpinner("Checking for update", "You are on the latest version!\n", "Checking for update failed, try again!\n", false)
		checkSpinner.Start()

		url := "http://localhost:8080/v1/dist/download/cli"
		params := map[string]string{"os": runtime.GOOS, "arch": runtime.GOARCH}

		updateUrl, shouldUpdate, checkErr := install.GetGitHubRelease(url, params, auxiliary.StateInstance.CLIVersion)
		if checkErr != nil {
			fmt.Println(checkErr)
		}

		CLIPath, pathErr := os.Executable()
		if pathErr != nil {
			fmt.Println(pathErr)
		}

		if shouldUpdate {
			if downloadErr := install.DownloadGitHubRelease("Pluralith CLI", updateUrl, CLIPath); downloadErr != nil {
				fmt.Println(downloadErr)
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
