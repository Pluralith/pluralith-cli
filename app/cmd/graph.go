package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"pluralith/pkg/auth"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/graph"
	"pluralith/pkg/install/components"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"
)

// stripCmd represents the strip command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Generate and export a diagram of the current plan state as a PDF",
	Long:  `Generate and export a diagram of the current plan state as a PDF`,
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

		// Check if graph module installed, if not -> install
		_, versionErr := exec.Command(filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing"), "version").Output()
		if versionErr != nil {
			components.GraphModule()
		}

		// Parse flag values
		diagramValues, valueErr := graph.GetDiagramValues(cmd.Flags())
		if valueErr != nil {
			fmt.Println(fmt.Errorf("getting diagram values failed -> %w", valueErr))
			return
		}

		// Run terraform plan to create execution plan if not specified otherwise by user
		if diagramValues["SkipPlan"] == false {
			_, planErr := terraform.RunPlan("plan", true)
			if planErr != nil {
				fmt.Println(fmt.Errorf("running terraform plan failed -> %w", planErr))
				return
			}
		} else {
			ux.PrintFormatted("→ ", []string{"bold", "blue"})
			ux.PrintFormatted("Plan\n", []string{"bold", "white"})
			ux.PrintFormatted("  -", []string{"blue", "bold"})
			fmt.Println(" Skipped\n")
		}

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

		// Pass plan state on to graphing module
		diagramValues["PlanStatePath"] = planStatePath

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
	graphCmd.PersistentFlags().Bool("generate-md", false, "Generate markdown output with exported PDF link and preview image for pull request comments")
}
