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
	"fmt"
	"pluralith/pkg/install/components"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

// stripCmd represents the strip command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Strip a given state file of secrets according to config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()

		ux.PrintFormatted("⠿ ", []string{"blue"})
		fmt.Println("Pass a component to install:\n")

		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Println(" graph-module")
		ux.PrintFormatted("→", []string{"blue", "bold"})
		fmt.Println(" ui\n")
	},
}

// Graph module
var installGraphModule = &cobra.Command{
	Use:   "graph-module",
	Short: "Strip a given state file of secrets according to config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		components.GraphModule()
	},
}

func init() {
	// installCmd.AddCommand(installGraphModule)
	rootCmd.AddCommand(installCmd)
}
