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
	Long:  `Set credentials for communication with the Pluralith API`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		fmt.Print("Welcome to ")
		ux.PrintFormatted("Pluralith!\n\n", []string{"blue"})

		APIKey, flagError := cmd.Flags().GetString("api-key")
		if flagError != nil {
			fmt.Println(fmt.Errorf("reading flag failed -> %w", flagError))
			return
		}

		// If no API key given via flag -> Prompt user for input
		if APIKey == "" {
			ux.PrintFormatted("→", []string{"blue", "bold"})
			fmt.Print("Enter API Key (You can find it in the Dashboard user settings https://app.pluralith.com/#/user/settings): ")

			// Capture user input
			fmt.Scanln(&APIKey)
		}

		if _, loginErr := auth.RunLogin(APIKey); loginErr != nil {
			fmt.Println(fmt.Errorf("failed to authenticate you -> %w", loginErr))
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().String("api-key", "", "Your Pluralith API key passed directly, to skip user prompt (for automation)")
}
