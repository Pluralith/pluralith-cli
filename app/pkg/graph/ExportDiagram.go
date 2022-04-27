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

func ExportDiagram(exportArgs map[string]interface{}) error {
	functionName := "ExportDiagram"

	ux.PrintFormatted("\n→", []string{"blue", "bold"})
	ux.PrintFormatted(" Export\n", []string{"white", "bold"})

	graphModulePath := filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing")

	exportSpinner := ux.NewSpinner("Generating Diagram PDF", "PDF Generated", "PDF Generation Failed", true)
	exportSpinner.Start()

	cmd := exec.Command(
		graphModulePath,
		"graph",
		"--apiKey", auxiliary.StateInstance.APIKey,
		"--title", exportArgs["Title"].(string),
		"--author", exportArgs["Author"].(string),
		"--ver", exportArgs["Version"].(string),
		"--fileName", exportArgs["FileName"].(string),
		"--outDir", exportArgs["OutDir"].(string),
		"--planStatePath", exportArgs["PlanStatePath"].(string),
		"--showChanges", strconv.FormatBool(exportArgs["ShowChanges"].(bool)),
		"--showDrift", strconv.FormatBool(exportArgs["ShowDrift"].(bool)),
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

	exportPath := filepath.Join(exportArgs["OutDir"].(string), exportArgs["FileName"].(string)) + ".pdf"

	exportSpinner.Success()
	ux.PrintFormatted("  → ", []string{"blue"})
	fmt.Print("Diagram Exported To: ")
	ux.PrintFormatted(exportPath+"\n", []string{"blue"})

	return nil
}
