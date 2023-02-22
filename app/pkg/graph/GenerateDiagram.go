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

func GenerateDiagram(exportArgs map[string]interface{}, costArgs map[string]interface{}, localRun bool) error {
	functionName := "GenerateDiagram"

	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Graph\n", []string{"white", "bold"})

	graphModulePath := filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing")

	var exportSpinner = ux.NewSpinner("Generating Diagram", "Diagram Generated", "Diagram Generation Failed", true)
	if localRun {
		exportSpinner = ux.NewSpinner("Generating Local Diagram", "Local Diagram Generated", "Local Diagram Generation Failed", true)
	}
	exportSpinner.Start()

	// Handle out dir and file name defaults
	if exportArgs["file-name"] == "" {
		exportArgs["file-name"] = strings.ReplaceAll(exportArgs["title"].(string), " ", "_")
	}
	if exportArgs["out-dir"] == "" {
		exportArgs["out-dir"] = auxiliary.StateInstance.WorkingPath
	}

	commandArgs := []string{
		"graph",
		"--api-key", auxiliary.StateInstance.APIKey,
		"--title", exportArgs["title"].(string),
		"--branch", exportArgs["branch"].(string),
		"--author", exportArgs["author"].(string),
		"--ver", exportArgs["version"].(string),
		"--file-name", exportArgs["file-name"].(string),
		"--out-dir", exportArgs["out-dir"].(string),
		"--plan-json-path", exportArgs["plan-json-path"].(string),
		"--cost-json-path", exportArgs["cost-json-path"].(string),
		"--export-pdf", strconv.FormatBool(exportArgs["export-pdf"].(bool)),
		"--sync-to-backend", strconv.FormatBool(exportArgs["sync-to-backend"].(bool)),
		"--show-changes", strconv.FormatBool(exportArgs["show-changes"].(bool)),
		"--show-drift", strconv.FormatBool(exportArgs["show-drift"].(bool)),
		"--show-costs", strconv.FormatBool(exportArgs["show-costs"].(bool)),
		"--wsl", strconv.FormatBool(auxiliary.StateInstance.IsWSL),

		"--cost-mode", costArgs["cost-mode"].(string),
		"--cost-period", costArgs["cost-period"].(string),
	}

	if auxiliary.StateInstance.PluralithConfig.Diagram != nil {
		commandArgs = append(commandArgs, "--config-json-path", filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.diagram.json"))
	}

	// If a run id is given (only on 'pluralith run') -> pass url components to graph module
	if exportArgs["runId"] != nil {
		commandArgs = append(
			commandArgs,
			"--run-id", exportArgs["runId"].(string),
			"--project-id", exportArgs["projectId"].(string),
			"--org-id", exportArgs["orgId"].(string),
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

	if exportArgs["export-pdf"] == true {
		ux.PrintFormatted("  → ", []string{"blue"})
		fmt.Print("Diagram Exported To: ")
		ux.PrintFormatted(exportPath+"\n", []string{"blue"})
	}

	return nil
}
