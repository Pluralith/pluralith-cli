package graph

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/install/components"
	"pluralith/pkg/terraform"
	"pluralith/pkg/ux"
)

func GenerateGraph(tfArgs map[string]interface{}, costArgs map[string]interface{}, exportArgs map[string]interface{}, runAsCI bool) error {
	functionName := "GenerateGraph"

	// Check if graph module installed, if not -> install
	_, versionErr := exec.Command(filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing"), "version").Output()
	if versionErr != nil {
		components.GraphModule()
	}

	_, planErr := terraform.RunPlan("plan", tfArgs, costArgs, true)
	if planErr != nil {
		return fmt.Errorf("running terraform plan failed -> %v: %w", functionName, planErr)
	}

	// Construct plan state path
	planJsonPath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.state.json")
	costJsonPath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.costs.json")

	// Check if plan state exists
	_, existErr := os.Stat(planJsonPath)     // Check if old state exists
	if errors.Is(existErr, os.ErrNotExist) { // If it exists -> delete
		ux.PrintFormatted("✘", []string{"bold", "red"})
		fmt.Print(" No plan state found ")
		return nil
	}

	// Pass plan state on to graphing module
	exportArgs["plan-json-path"] = planJsonPath
	exportArgs["cost-json-path"] = costJsonPath

	// Generate diagram through graphing module
	if diagramErr := GenerateDiagram(exportArgs, costArgs); diagramErr != nil {
		return fmt.Errorf("generating diagram failed -> %v: %w", functionName, diagramErr)
	}

	return nil
}
