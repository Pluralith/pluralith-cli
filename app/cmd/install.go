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
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/install"
	"pluralith/pkg/ux"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Strip a given state file of secrets according to config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		ux.PrintFormatted("⠿ ", []string{"blue"})
		fmt.Println("Pass a component to install:\n")

		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Println(" graph-module")
		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Println(" ui\n")
	},
}

var installGraphModule = &cobra.Command{
	Use:   "graph-module",
	Short: "Strip a given state file of secrets according to config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		fmt.Print("Installing Latest ")
		ux.PrintFormatted("Graph Module\n\n", []string{"bold", "blue"})

		// Construct url
		url := "http://localhost:8080/v1/dist/download/cli/graphing"
		params := map[string]string{"os": runtime.GOOS, "arch": runtime.GOARCH}

		// Generate install path
		installPath := filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing")

		// Get current version
		var currentVersion string
		currentVersionByte, versionErr := exec.Command(installPath, "version").Output()
		if versionErr != nil {
			currentVersion = ""
		} else {
			currentVersion = strings.TrimSpace(string(currentVersionByte))
		}

		// Get Github release
		downloadUrl, shouldDownload, checkErr := install.GetGitHubRelease(url, params, currentVersion)
		if checkErr != nil {
			fmt.Println(checkErr)
		}

		// Handle download
		if shouldDownload {
			if downloadErr := install.DownloadGitHubRelease("Graph Module", downloadUrl, installPath); downloadErr != nil {
				fmt.Println(downloadErr)
			}
		}
	},
}

func init() {
	installCmd.AddCommand(installGraphModule)
	rootCmd.AddCommand(installCmd)
}
