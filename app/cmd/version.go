package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
	"strings"

	"github.com/spf13/cobra"
)

// planOldCmd represents the planOld command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "List the installed versions of the Pluralith CLI and its modules",
	Long:  `List the installed versions of the Pluralith CLI and its modules`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		// Get graph module version
		var graphVersion string
		currentVersionByte, versionErr := exec.Command(filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing"), "version").Output()
		if versionErr != nil {
			graphVersion = "Not Installed"
		} else {
			graphVersion = strings.TrimSpace(string(currentVersionByte))
		}

		ux.PrintFormatted("→ ", []string{"blue"})
		fmt.Print("CLI Version: ")
		ux.PrintFormatted(auxiliary.StateInstance.CLIVersion+"\n", []string{"bold", "blue"})

		ux.PrintFormatted("→ ", []string{"blue"})
		fmt.Print("Graph Module Version: ")
		ux.PrintFormatted(graphVersion+"\n\n", []string{"bold", "blue"})
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
