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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/graph"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"
)

// stripCmd represents the strip command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Generate and export diagram of the current plan state as a PDF",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Verify API key with backend
		isValid, verifyErr := auth.VerifyAPIKey(auxiliary.StateInstance.APIKey)
		if verifyErr != nil {
			fmt.Println(fmt.Errorf("verifying API key failed -> %w", verifyErr))
			return
		}
		if !isValid {
			ux.PrintFormatted("\n✖️", []string{"red", "bold"})
			fmt.Print(" Invalid API key → Run ")
			ux.PrintFormatted("pluralith login", []string{"blue"})
			fmt.Println(" again\n")
			return
		}

		// Parse flag values
		diagramValues, valueErr := graph.GetDiagramValues(cmd.Flags())
		if valueErr != nil {
			fmt.Println(fmt.Errorf("getting diagram values failed -> %w", valueErr))
			return
		}

		// Run terraform plan to create execution plan if not specified otherwise by user
		if diagramValues["SkipPlan"] == false {
			_, planErr := terraform.RunPlan("plan")
			if planErr != nil {
				fmt.Println(fmt.Errorf("running terraform plan failed -> %w", planErr))
				return
			}
		} else {
			ux.PrintFormatted("→ ", []string{"bold", "blue"})
			ux.PrintFormatted("Plan\n", []string{"bold", "white"})
			ux.PrintFormatted("  -", []string{"blue", "bold"})
			fmt.Println(" Skipped")
		}

		fmt.Println()

		// Construct plan state path
		planStatePath := filepath.Join(auxiliary.StateInstance.WorkingPath, "pluralith.state.stripped")

		// Check if plan state exists
		_, existErr := os.Stat(planStatePath)    // Check if old state exists
		if errors.Is(existErr, os.ErrNotExist) { // If it exists -> delete
			ux.PrintFormatted("✖️", []string{"bold", "red"})
			fmt.Print(" No plan state found ")
			ux.PrintFormatted("→", []string{"bold", "red"})
			fmt.Println(" Run pluralith graph again without --skip-plan")
			return
		}

		// Read plan state json
		planStateString, readErr := os.ReadFile(planStatePath)
		if readErr != nil {
			fmt.Println(fmt.Errorf("running terraform plan failed -> %w", readErr))
			return
		}

		// Pass plan state on to graphing module
		diagramValues["PlanState"] = string(planStateString)

		// Generate diagram through graphing module
		if exportErr := graph.ExportDiagram(diagramValues); exportErr != nil {
			fmt.Println(fmt.Errorf("exporting diagram failed -> %w", exportErr))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
	graphCmd.PersistentFlags().String("title", "", "The title for your diagram, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("author", "", "The author/creator of the diagram, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("version", "", "The diagram version, will be displayed in the PDF output")
	graphCmd.PersistentFlags().String("out-dir", "", "The directory the diagram should be exported to")
	graphCmd.PersistentFlags().String("file-name", "", "The name of the exported PDF")
	graphCmd.PersistentFlags().Bool("skip-plan", false, "Generates a diagram without running plan again (needs pluralith state from previous plan run)")
}
