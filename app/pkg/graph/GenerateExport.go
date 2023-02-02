package graph

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
	"strconv"
	"strings"
)

func GenerateExport(exportArgs map[string]interface{}, costArgs map[string]interface{}) error {
	functionName := "GenerateExport"

	graphModulePath := filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing")

	fmt.Println()
	var exportSpinner = ux.NewSpinner("Generating Export", "Export Generated", "Export Generation Failed", true)
	exportSpinner.Start()

	// Handle out dir and file name defaults
	if exportArgs["file-name"] == "" {
		exportArgs["file-name"] = strings.ReplaceAll(exportArgs["title"].(string), " ", "_")
	}
	if exportArgs["out-dir"] == "" {
		exportArgs["out-dir"] = auxiliary.StateInstance.WorkingPath
	}

	commandArgs := []string{
		"export",
		"--file-name", exportArgs["file-name"].(string),
		"--title", exportArgs["title"].(string),
		"--branch", exportArgs["branch"].(string),
		"--author", exportArgs["author"].(string),
		"--ver", exportArgs["version"].(string),
		"--out-dir", exportArgs["out-dir"].(string),
		"--sync-to-backend", strconv.FormatBool(exportArgs["sync-to-backend"].(bool)),
		"--show-changes", strconv.FormatBool(exportArgs["show-changes"].(bool)),
		"--show-drift", strconv.FormatBool(exportArgs["show-drift"].(bool)),
		"--show-costs", strconv.FormatBool(exportArgs["show-costs"].(bool)),
		"--cost-mode", costArgs["cost-mode"].(string),
		"--cost-period", costArgs["cost-period"].(string),
	}

	// If a run id is given (only on 'pluralith run') -> pass url components to graph module
	if exportArgs["run-id"] != nil {
		commandArgs = append(
			commandArgs,
			"--run-id", exportArgs["run-id"].(string),
			"--project-id", exportArgs["project-id"].(string),
			"--org-id", exportArgs["org-id"].(string),
		)
	}

	cmd := exec.Command(
		graphModulePath,
		commandArgs...,
	)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink
	cmd.Stdin = os.Stdin

	if runErr := cmd.Run(); runErr != nil {
		exportSpinner.Fail()
		fmt.Println(errorSink.String())
		return fmt.Errorf("running CLI command failed -> %v: %w", functionName, runErr)
	}

	exportPath := filepath.Join(exportArgs["out-dir"].(string), exportArgs["file-name"].(string)) + ".pdf"

	exportSpinner.Success()

	ux.PrintFormatted("  → ", []string{"blue"})
	fmt.Print("Diagram Exported To: ")
	ux.PrintFormatted(exportPath+"\n", []string{"blue"})

	return nil
}
