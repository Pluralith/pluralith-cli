package cmd

import (
	"fmt"
	"pluralith/cmd/run"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Generate a diagram in CI and post it to the Pluralith dashboard",
	Long:  `Generate a diagram in CI and post it to the Pluralith dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		ux.PrintFormatted("⠿ ", []string{"blue"})
		fmt.Print("Append the Terraform command you'd like to run:\n\n")

		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Println(" pluralith run plan")
		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Print(" pluralith run apply\n\n")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.AddCommand(run.RunPlanCmd)
	runCmd.AddCommand(run.RunApplyCmd)
	runCmd.AddCommand(run.RunDestroyCmd)
	runCmd.PersistentFlags().String("title", "", "The title for your diagram, will be displayed in the PDF output")
	runCmd.PersistentFlags().String("version", "", "The diagram version, will be displayed in the PDF output")
	runCmd.PersistentFlags().String("out-dir", "", "The directory the diagram should be exported to")
	runCmd.PersistentFlags().String("file-name", "", "The name of the exported PDF")
	runCmd.PersistentFlags().Bool("post-apply", false, "Determines whether this run is after an apply and should update the latest docs for this infrastructure project")
	runCmd.PersistentFlags().Bool("export-pdf", false, "Determines whether a PDF export of the run Diagram is generated locally")
	runCmd.PersistentFlags().Bool("sync-to-backend", false, "Determines whether a PDF export of your diagram is stored in your state backend alongside your Terraform state. (Currently supports azurerm, gcs and s3)")
	runCmd.PersistentFlags().Bool("show-changes", false, "Determines whether the exported diagram highlights changes made in the latest Terraform plan or outputs a general diagram of the infrastructure")
	runCmd.PersistentFlags().Bool("show-drift", false, "Determines whether the exported diagram highlights resource drift detected by Terraform")
	runCmd.PersistentFlags().Bool("show-costs", false, "Determines whether the exported diagram includes cost information")
	runCmd.PersistentFlags().String("cost-mode", "delta", "Determines which costs are shown. Can be 'delta' or 'total'")
	runCmd.PersistentFlags().String("cost-period", "month", "Determines over which period costs are aggregated. Can be 'hour' or 'month'")
	runCmd.PersistentFlags().StringArray("var-file", []string{}, "Path to a var file to pass to Terraform. Can be specified multiple times.")
	runCmd.PersistentFlags().StringArray("var", []string{}, "A variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')")
	runCmd.PersistentFlags().String("plan-file", "", "Path to an execution plan binary file. If passed, this will skip a plan run under the hood.")
	runCmd.PersistentFlags().String("plan-file-json", "", "Path to an execution plan json file. If passed, this will skip a plan run under the hood.")
	runCmd.PersistentFlags().String("cost-usage-file", "", "Path to an infracost usage file to be used for the cost breakdown")
}
