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
	"pluralith/helpers"
	"pluralith/ux"

	"github.com/spf13/cobra"
)

// Defining command args/flags
var pluralithDestroyArgs = []string{}

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy and show changes in Pluralith",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Initializing variable for manual user confirmation
		var confirm string
		// Manually parsing args (due to cobra lacking a feature)
		parsedArgs, parsedArgMap := helpers.ParseArgs(args, pluralithDestroyArgs)

		// Checking if auto-approve flag has been set
		if parsedArgMap["auto-approve"] == "" {
			// Appending -auto-approve flag to run command properly after confirmation
			parsedArgs = append(parsedArgs, "-auto-approve")
			// Handling UX and user input
			ux.PrintFormatted("?", []string{"blue", "bold"})
			fmt.Println(" Destroy Current Infrastructure?")
			ux.PrintFormatted("  Yes to confirm: ", []string{"bold"})
			fmt.Scanln(&confirm)
		}

		// If user confirms manually or auto-approve flag has been set -> Run destroy
		if confirm == "yes" || parsedArgMap["auto-approve"] != "" {
			ux.PrintFormatted("\n✔", []string{"blue", "bold"})
			fmt.Println(" Destruction Confirmed")

			// Launching Pluralith
			helpers.LaunchPluralith()

			ux.PrintFormatted("⠿", []string{"blue", "bold"})
			fmt.Println(" Destruction Status:")

			// Running destroy command with args passed by user
			if destroyOutput, destroyErr := helpers.ExecuteTerraform("destroy", parsedArgs, true); destroyErr != nil {
				// Handling failed terraform destroy
				ux.PrintFormatted("✖️", []string{"red", "bold"})
				fmt.Println(" Destroy Failed")
				fmt.Println(destroyOutput)
			} else {
				// Handling successful terraform destroy
				ux.PrintFormatted("✔ All Done!\n", []string{"blue", "bold"})
			}

			// Updating command in hist to update Pluralith UI
			helpers.WriteToHist("destroy", "terraform-end\n")
		} else {
			ux.PrintFormatted("\n✖️", []string{"red", "bold"})
			fmt.Println(" Destroy Aborted")
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
