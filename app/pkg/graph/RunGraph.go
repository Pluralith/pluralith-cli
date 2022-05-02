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

func RunGraph(tfArgs []string, costArgs []string, exportArgs map[string]interface{}, runAsCI bool) error {
	functionName := "RunGraph"

	// Check if graph module installed, if not -> install
	_, versionErr := exec.Command(filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing"), "version").Output()
	if versionErr != nil {
		components.GraphModule()
	}

	// Run terraform plan to create execution plan if not specified otherwise by user
	if exportArgs["SkipPlan"] == false {
		_, planErr := terraform.RunPlan("plan", tfArgs, costArgs, true)
		if planErr != nil {
			return fmt.Errorf("running terraform plan failed -> %v: %w", functionName, planErr)
		}
	} else {
		ux.PrintFormatted("→ ", []string{"bold", "blue"})
		ux.PrintFormatted("Plan\n", []string{"bold", "white"})
		ux.PrintFormatted("  -", []string{"blue", "bold"})
		fmt.Println(" Skipped\n")
	}

	// Construct plan state path
	planStatePath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.state.json")

	// Check if plan state exists
	_, existErr := os.Stat(planStatePath)    // Check if old state exists
	if errors.Is(existErr, os.ErrNotExist) { // If it exists -> delete
		ux.PrintFormatted("✘", []string{"bold", "red"})
		fmt.Print(" No plan state found ")
		ux.PrintFormatted("→", []string{"bold", "red"})
		fmt.Println(" Run pluralith graph again without --skip-plan")
		return nil
	}

	// Pass plan state on to graphing module
	exportArgs["PlanStatePath"] = planStatePath

	// Generate diagram through graphing module
	if exportErr := ExportDiagram(exportArgs); exportErr != nil {
		return fmt.Errorf("exporting diagram failed -> %v: %w", functionName, exportErr)
	}

	if runAsCI {
		if exportErr := HandleCIRun(exportArgs); exportErr != nil {
			return fmt.Errorf("exporting diagram failed -> %v: %w", functionName, exportErr)
		}
	}

	fmt.Println()

	return nil
}
