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
	"io/ioutil"
	"pluralith/helpers"
	"pluralith/ux"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var stripCmd = &cobra.Command{
	Use:   "strip",
	Short: "Strip a given state file of secrets according to config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Fetching all state files in current working directory
		stateFiles := helpers.FetchFiles(".tfstate")

		// Instantiating new strip spinner
		stripSpinner := ux.NewSpinner("Stripping Secrets", fmt.Sprintf("Secrets Stripped From %d File", len(stateFiles)), "Stripping Secrets Failed")
		stripSpinner.Start()

		// Stripping secrets and writing stripped state to disk
		for fileName, fileContent := range stateFiles {
			strippedFile, err := helpers.StripSecrets(fileContent, sensitiveKeys, "gatewatch")
			if err != nil {
				stripSpinner.Fail("Failed to strip secrets from %s", fileName)
			} else {
				ioutil.WriteFile(fmt.Sprintf("%s.plstate.stripped", fileName), []byte(strippedFile), 0644)
				stripSpinner.Success()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(stripCmd)
}
