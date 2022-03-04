package cmd

import (
	"fmt"
	"pluralith/pkg/install/components"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install additional modules to expand the CLI's featureset",
	Long:  `Install additional modules to expand the CLI's featureset`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		ux.PrintFormatted("⠿ ", []string{"blue"})
		fmt.Println("Pass a component to install:\n")

		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Println(" graph-module")
		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Println(" ui (coming soon)\n")
	},
}

// Graph module
var installGraphModule = &cobra.Command{
	Use:   "graph-module",
	Short: "Install the latest graph module for the Pluralith CLI",
	Long:  `Install the latest graph module for the Pluralith CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		components.GraphModule()
	},
}

func init() {
	installCmd.AddCommand(installGraphModule)
	rootCmd.AddCommand(installCmd)
}
