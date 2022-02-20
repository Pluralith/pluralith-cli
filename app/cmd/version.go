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
	"pluralith/pkg/ux"
	"strings"

	"github.com/spf13/cobra"
)

// planOldCmd represents the planOld command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Check for and install updates for the Pluralith CLI",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		// Get graph module version
		var graphVersion string
		currentVersionByte, versionErr := exec.Command(filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing"), "version").Output()
		if versionErr != nil {
			graphVersion = "Not Installed"
		} else {
			graphVersion = strings.TrimSpace(string(currentVersionByte))
		}

		ux.PrintFormatted("→ ", []string{"blue"})
		fmt.Print("CLI Version: ")
		ux.PrintFormatted(auxiliary.StateInstance.CLIVersion+"\n", []string{"bold", "blue"})

		ux.PrintFormatted("→ ", []string{"blue"})
		fmt.Print("Graph Module Version: ")
		ux.PrintFormatted(graphVersion+"\n\n", []string{"bold", "blue"})
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
