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
	"net/http"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
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

		verificationSpinner := ux.NewSpinner("Verifying your API key", "Your API key is valid, you are logged in!\n", "API key verification failed\n", false)

		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Print(" Enter API Key: ")

		// Capture user input
		var APIKey string
		fmt.Scanln(&APIKey)

		verificationSpinner.Start()

		// Construct key verification request
		request, _ := http.NewRequest("GET", "http://localhost:8080/v1/auth/key/verify", nil)
		request.Header.Add("Authorization", "Bearer "+APIKey)

		// Execute key verification request
		client := &http.Client{}
		response, responseErr := client.Do(request)

		if responseErr != nil {
			verificationSpinner.Fail("Failed to verify API key\n")
			fmt.Println(fmt.Errorf("%w", responseErr))
		}

		// Hande verification response
		if response.StatusCode == 200 {
			auxiliary.StateInstance.APIKey = APIKey
			credentialsPath := filepath.Join(auxiliary.StateInstance.PluralithPath, "credentials")

			// Write api key to credentials file
			if writeErr := os.WriteFile(credentialsPath, []byte(APIKey), 0700); writeErr != nil {
				verificationSpinner.Fail("Failed to write API key to config try again!\n")
				fmt.Println(fmt.Errorf("%w", writeErr))
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
