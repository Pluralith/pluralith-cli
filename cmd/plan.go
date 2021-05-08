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
	"pluralith/helpers"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var pluralithArgs = []string{"-show-output", "-s"}

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan and draw diagram",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Defining blue print function
		printBlue := color.New(color.FgBlue, color.Bold).PrintfFunc()

		// Manually parsing arg (due to cobra lacking a feature)
		parsedArgs, parsedArgMap := helpers.ParseArgs(args, pluralithArgs)
		// Getting value of -out flag
		planOut := parsedArgMap["out"]

		// If no value is given for -out, replace it with standard ./pluralith
		if planOut == "" {
			planOut = "./pluralith"
			parsedArgs = append(parsedArgs, "-out", planOut)
		}

		// Running terraform plan command with cleaned up args to generate execution plan
		if _, code := helpers.ExecuteTerraform("plan", parsedArgs, false, false); code == 0 {
			// If plan command succeeds -> Run terraform show on previously generated execution plan to generate plan state file
			helpers.ExecuteTerraform("show", []string{"-json", planOut}, false, false)
			printBlue("\n✔ All Done!\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
