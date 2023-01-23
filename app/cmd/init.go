package cmd

import (
	"fmt"
	"pluralith/pkg/initialization"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// planOldCmd represents the planOld command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Pluralith project in the current directory",
	Long:  `Initialize a Pluralith project in the current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		fmt.Print("Welcome to ")
		ux.PrintFormatted("Pluralith!\n", []string{"blue"})

		// Get flag values
		initData := initialization.InitData{}
		isEmpty, _ := cmd.Flags().GetBool("empty")

		if isEmpty {
			if writeErr := initialization.WriteConfig(initData); writeErr != nil {
				fmt.Println(fmt.Errorf("failed to create config template -> %w", writeErr))
			}
			return
		}

		initData.APIKey, _ = cmd.Flags().GetString("api-key")
		initData.OrgId, _ = cmd.Flags().GetString("org-id")
		initData.ProjectId, _ = cmd.Flags().GetString("project-id")
		initData.ProjectName, _ = cmd.Flags().GetString("project-name")

		noInputs, _ := cmd.Flags().GetBool("no-inputs")

		_, _, initErr := initialization.RunInit(noInputs, initData, false)
		if initErr != nil {
			fmt.Println(fmt.Errorf("pluralith init failed -> %w", initErr))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().String("api-key", "", "Your Pluralith API key. Pass via flag to skip user prompt.")
	initCmd.PersistentFlags().String("org-id", "", "Your Org Id (Can be found in your Pluralith dashboard). Pass via flag to skip user prompt and override pluralith.yml")
	initCmd.PersistentFlags().String("project-id", "", "Your Project Id (If no project with passed Id exists, one gets created). Pass via flag to skip user prompt and override pluralith.yml")
	initCmd.PersistentFlags().String("project-name", "", "Your Project name. Pass via flag to skip user prompt and override pluralith.yml")
	initCmd.PersistentFlags().Bool("empty", false, "Creates an empty pluralith.yml config file in the current directory")
	initCmd.PersistentFlags().Bool("no-inputs", false, "Run init assuming all necessary parameters are present (for CI)")
}
