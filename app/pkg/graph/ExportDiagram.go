package graph

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func ExportDiagram(diagramValues map[string]string) error {
	functionName := "ExportDiagram"

	graphModulePath := filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing")

	exportSpinner := ux.NewSpinner("Generating Diagram PDF", "PDF Export successful!", "PDF Export failed", false)
	exportSpinner.Start()

	cmd := exec.Command(
		graphModulePath,
		"graph",
		"--apiKey", auxiliary.StateInstance.APIKey,
		"--title", diagramValues["Title"],
		"--author", diagramValues["Author"],
		"--ver", diagramValues["Version"],
		"--fileName", diagramValues["FileName"],
		"--outDir", diagramValues["OutDir"],
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

	exportSpinner.Success()
	ux.PrintFormatted("→ ", []string{"blue"})
	fmt.Print("Diagram exported to: ")
	ux.PrintFormatted(filepath.Join(diagramValues["OutDir"], diagramValues["FileName"])+".pdf", []string{"blue"})
	fmt.Println("\n")

	return nil
}
