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
package main

import (
	"fmt"
	"pluralith/cmd"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/dblock"
)

// Initialize various components of application
func initApp() {
	functionName := "initApp"

	dblock.LockInstance.GenerateLock()
	auxiliary.StateInstance.CLIVersion = "0.1.3"
	fmt.Println("update-test-build-v2")

	if pathGenErr := auxiliary.StateInstance.GeneratePaths(); pathGenErr != nil {
		fmt.Println(fmt.Errorf("generating application paths failed -> %v: %w", functionName, pathGenErr))
	}
	if pathInitErr := auxiliary.StateInstance.InitPaths(); pathInitErr != nil {
		fmt.Println(fmt.Errorf("initializing application directories failed -> %v: %w", functionName, pathInitErr))
	}
	if setAPIKeyErr := auxiliary.StateInstance.SetAPIKey(); setAPIKeyErr != nil {
		fmt.Println(fmt.Errorf("setting API key failed -> %v: %w", functionName, setAPIKeyErr))
	}

	auxiliary.StateInstance.CheckCI()

	if filterInitErr := auxiliary.FilterInstance.InitFilters(); filterInitErr != nil {
		fmt.Println(fmt.Errorf("initializing secret filters failed -> %v: %w", functionName, filterInitErr))
	}
	if getConfigErr := auxiliary.FilterInstance.GetSecretConfig(); getConfigErr != nil {
		fmt.Println(fmt.Errorf("fetching secret config failed -> %v: %w", functionName, getConfigErr))
	}
}

func main() {
	initApp()
	cmd.Execute()
}
