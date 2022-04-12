package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auth"
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
		fmt.Println("Lets set up your project and get you up and running.\n")

		var APIKey string
		var ProjectId int

		if auxiliary.StateInstance.APIKey == "" {
			ux.PrintFormatted("⠿ ", []string{"blue"})
			fmt.Println("We noticed you are not authenticated!")
			ux.PrintFormatted("→", []string{"blue", "bold"})
			fmt.Print(" Enter your API Key: ")

			// Capture user input
			// var APIKey string
			fmt.Scanln(&APIKey)
			loginValid, loginErr := auth.RunLogin(APIKey)
			if !loginValid {
				return
			}
			if loginErr != nil {
				fmt.Println(loginErr)
			}
		}

		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Print(" Enter Project Id: ")

		// Capture user input
		fmt.Scanln(&ProjectId)

		configPath := filepath.Join(auxiliary.StateInstance.WorkingPath, "pluralith.yml")
		configString := fmt.Sprintf(initialization.ConfigTemplate, ProjectId)

		helperWriteErr := os.WriteFile(configPath, []byte(configString), 0700)
		if helperWriteErr != nil {
			fmt.Println(fmt.Errorf("failed to create config template -> %w", helperWriteErr))
			return
		}

		ux.PrintFormatted("  ✔", []string{"blue", "bold"})
		fmt.Print(" Your project has been initialized! Customize your config in ")
		ux.PrintFormatted("pluralith.yml\n\n", []string{"blue"})

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
