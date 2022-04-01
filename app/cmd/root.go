package cmd

import (
	"fmt"
	"os"

	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
)

var cfgFile string

// Defining custom "Long" message with colored output
var blueColor = color.New(color.FgHiBlue).SprintFunc()
var longText = fmt.Sprintf(`%s

Welcome to %s, a tool to visualize your Terraform state.
It hooks directly into your Terraform installation and draws up a realtime infrastructure diagram for you.
	
We are currently in early private alpha.
Give the project a closer look at:
%s`,
	blueColor(` _
|_)|    _ _ |._|_|_
|  ||_|| (_||| | | |`),
	blueColor("Pluralith"),
	blueColor("https://www.pluralith.com"),
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pluralith",
	Short: "An application for Terraform state visualisation",
	Long:  longText,
	Run: func(cmd *cobra.Command, args []string) {
		ux.PrintHead()
		auxiliary.LaunchPluralith()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pluralith.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pluralith" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pluralith")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
