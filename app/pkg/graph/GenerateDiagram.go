package graph

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ci"
	"pluralith/pkg/ux"
	"strconv"
	"strings"
)

func GenerateComment(runCache map[string]interface{}) error {
	functionName := "preparePRComment"

	// Generate pull request comment markdown
	commentMD, commentErr := ci.GenerateMD(runCache["urls"].(map[string]interface{}), runCache["changes"].(map[string]interface{}))
	if commentErr != nil {
		return fmt.Errorf("generating PR comment markdown failed -> %v: %w", functionName, commentErr)
	}

	// Write markdown to file system for usage by pipeline
	commentPath := filepath.Join(auxiliary.StateInstance.WorkingPath, "comment.md")
	if writeErr := os.WriteFile(commentPath, []byte(commentMD), 0700); writeErr != nil {
		return fmt.Errorf("writing PR comment markdown to filesystem failed -> %v: %w", functionName, writeErr)
	}

	return nil
}

func GenerateDiagram(exportArgs map[string]interface{}) error {
	functionName := "GenerateDiagram"

	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Export\n", []string{"white", "bold"})

	graphModulePath := filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing")

	exportSpinner := ux.NewSpinner("Generating Diagram", "Diagram Generated", "Diagram Generation Failed", true)
	exportSpinner.Start()

	// Handle out dir and file name defaults
	if exportArgs["file-name"] == "" {
		exportArgs["file-name"] = strings.ReplaceAll(exportArgs["title"].(string), " ", "_")
	}
	if exportArgs["out-dir"] == "" {
		exportArgs["out-dir"] = auxiliary.StateInstance.WorkingPath
	}

	cmd := exec.Command(
		graphModulePath,
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
		"--show-changes", strconv.FormatBool(exportArgs["show-changes"].(bool)),
		"--show-drift", strconv.FormatBool(exportArgs["show-drift"].(bool)),
		"--show-costs", strconv.FormatBool(exportArgs["show-costs"].(bool)),
		"--export-pdf", strconv.FormatBool(exportArgs["export-pdf"].(bool)),
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
