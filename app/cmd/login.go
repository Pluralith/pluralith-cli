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
			fmt.Print(" Enter API Key: ")

			// Capture user input
			// var APIKey string
			fmt.Scanln(&APIKey)
		}

		verificationSpinner := ux.NewSpinner("Verifying your API key", "Your API key is valid, you are logged in!\n", "API key verification failed\n", false)
		verificationSpinner.Start()

		// Verify API key with backend
		isValid, verifyErr := auth.VerifyAPIKey(APIKey)
		if verifyErr != nil {
			fmt.Println(fmt.Errorf("verifying API key failed -> %w", verifyErr))
			return
		}

		if isValid {
			// Set API key in credentials file at ~/Pluralith/credentials
			setErr := auth.SetAPIKey(APIKey)
			if setErr != nil {
				verificationSpinner.Fail("Could not write to credentials file\n")
				fmt.Println(fmt.Errorf("setting API key in credentials file failed -> %w", setErr))
				return
			}
			verificationSpinner.Success()
		} else {
			verificationSpinner.Fail("The passed API key is invalid, try again!\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().String("api-key", "", "The Pluralith API key passed directly, skips user prompt (for automation)")
}
