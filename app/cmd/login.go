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
	"pluralith/pkg/auth"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Set credentials for communication with the Pluralith API",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		fmt.Print("Welcome to ")
		ux.PrintFormatted("Pluralith!\n\n", []string{"blue"})

		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Print(" Enter API Key: ")

		// Capture user input
		var APIKey string
		fmt.Scanln(&APIKey)

		verificationSpinner := ux.NewSpinner("Verifying your API key", "Your API key is valid, you are logged in!\n", "API key verification failed\n", false)
		verificationSpinner.Start()

		// Verify API key with backend
		isValid, verifyErr := auth.VerifyAPIKey(APIKey)
		if verifyErr != nil {
			fmt.Println(fmt.Errorf("verifying API key failed -> %w", verifyErr))
		}

		if isValid {
			// Set API key in credentials file at ~/Pluralith/credentials
			setErr := auth.SetAPIKey(APIKey)
			if setErr != nil {
				verificationSpinner.Fail("Could not write to credentials file\n")
				fmt.Println(fmt.Errorf("setting API key in credentials file failed -> %w", setErr))
			}
			verificationSpinner.Success()
		} else {
			verificationSpinner.Fail("The passed API key is invalid, try again!\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
