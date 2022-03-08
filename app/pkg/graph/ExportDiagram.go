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
)

func HandlePRUpdate(diagramPath string) error {
	functionName := "preparePRComment"

	// Upload diagram to storage for pull request comment hosting
	urls, hostErr := HostExport(diagramPath)
	if hostErr != nil {
		return fmt.Errorf("hosting diagram for PR comment failed -> %v: %w", functionName, hostErr)
	}

	// Generate pull request comment markdown
	commentMD, commentErr := ci.GeneratePRComment(urls)
	if commentErr != nil {
		return fmt.Errorf("hosting diagram for PR comment failed -> %v: %w", functionName, hostErr)
	}

	fmt.Println(commentMD)

	// Write markdown to file system for usage by pipeline
	// if writeErr := os.WriteFile("commend.md", []byte(commentMD), 0700); writeErr != nil {
	// 	return fmt.Errorf("writing PR comment markdown to filesystem failed -> %v: %w", functionName, hostErr)
	// }

	return nil
}

func ExportDiagram(diagramValues map[string]interface{}) error {
	functionName := "ExportDiagram"

	ux.PrintFormatted("→", []string{"blue", "bold"})
	ux.PrintFormatted(" Export\n", []string{"white", "bold"})

	graphModulePath := filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing")

	exportSpinner := ux.NewSpinner("Generating Diagram PDF", "PDF Export Successful!", "PDF Export Failed", true)
	exportSpinner.Start()

	cmd := exec.Command(
		graphModulePath,
		"graph",
		"--apiKey", auxiliary.StateInstance.APIKey,
		"--title", diagramValues["Title"].(string),
		"--author", diagramValues["Author"].(string),
		"--ver", diagramValues["Version"].(string),
		"--fileName", diagramValues["FileName"].(string),
		"--outDir", diagramValues["OutDir"].(string),
		"--planStatePath", diagramValues["PlanStatePath"].(string),
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

	exportPath := filepath.Join(diagramValues["OutDir"].(string), diagramValues["FileName"].(string)) + ".pdf"

	if auxiliary.StateInstance.IsCI {
		if prErr := HandlePRUpdate(exportPath); prErr != nil {
			return fmt.Errorf("handling pull request update failed -> %v: %w", functionName, prErr)
		}
	} else {
		if logErr := LogExport(); logErr != nil {
			return fmt.Errorf("logging diagram export failed -> %v: %w", functionName, logErr)
		}
	}

	exportSpinner.Success()
	ux.PrintFormatted("  → ", []string{"blue"})
	fmt.Print("Diagram exported to: ")
	ux.PrintFormatted(exportPath, []string{"blue"})
	fmt.Println("\n")

	return nil
}
