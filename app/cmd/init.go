package cmd

import (
	"fmt"
	"pluralith/pkg/auxiliary"
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

		APIKey := ""
		projectId := ""

		// Get flag values
		isEmpty, emptyError := cmd.Flags().GetBool("empty")
		if emptyError != nil {
			fmt.Println(fmt.Errorf("reading flag failed -> %w", emptyError))
			return
		}

		APIKey, APIKeyError := cmd.Flags().GetString("api-key")
		if APIKeyError != nil {
			fmt.Println(fmt.Errorf("reading flag failed -> %w", APIKeyError))
			return
		}

		// If no API key is passed, set to existing API key value in state (can be "" as well)
		if APIKey == "" {
			APIKey = auxiliary.StateInstance.APIKey
		}

		projectId, projectIdErr := cmd.Flags().GetString("project-id")
		if APIKeyError != nil {
			fmt.Println(fmt.Errorf("reading flag failed -> %w", projectIdErr))
			return
		}

		if initErr := initialization.RunInit(isEmpty, APIKey, projectId); initErr != nil {
			fmt.Println(fmt.Errorf("pluralith init failed -> %w", initErr))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().String("api-key", "", "Your Pluralith API key passed directly, to skip user prompt (for automation)")
	initCmd.PersistentFlags().String("project-id", "", "Your project id passed directly, to skip user prompt (for automation)")
	initCmd.PersistentFlags().Bool("empty", false, "Creates an empty pluralith.yml config file in the current directory")
}
