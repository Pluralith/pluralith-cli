/*
Copyright © 2022 Pluralith Industries Inc. founders@pluralith.com
*/
package main

import (
	"fmt"
	"pluralith/cmd"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/dblock"
	"pluralith/pkg/install/components"
)

// Initialize various components of application
func initApp() {
	functionName := "initApp"

	dblock.LockInstance.GenerateLock()
	auxiliary.StateInstance.CLIVersion = "0.1.17"

	if pathGenErr := auxiliary.StateInstance.GeneratePaths(); pathGenErr != nil {
		fmt.Println(fmt.Errorf("generating application paths failed -> %v: %w", functionName, pathGenErr))
	}
	if pathInitErr := auxiliary.StateInstance.InitPaths(); pathInitErr != nil {
		fmt.Println(fmt.Errorf("initializing application directories failed -> %v: %w", functionName, pathInitErr))
	}
	if setAPIKeyErr := auxiliary.StateInstance.SetAPIKey(); setAPIKeyErr != nil {
		fmt.Println(fmt.Errorf("setting API key failed -> %v: %w", functionName, setAPIKeyErr))
	}

	// Check for and install potential graph module update
	components.GraphModule(true)

	auxiliary.StateInstance.CheckCI()
	auxiliary.StateInstance.GetBranch()
	auxiliary.StateInstance.CheckTerraformInit()
	auxiliary.StateInstance.CheckInfracost()

	if getConfigErr := auxiliary.StateInstance.GetConfig(); getConfigErr != nil {
		fmt.Println(fmt.Errorf("fetching pluralith config failed -> %v: %w", functionName, getConfigErr))
	}
}

func main() {
	initApp()
	cmd.Execute()
}
