package graph

import (
	"bufio"
	"bytes"
	"encoding/json"
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

func GenerateComment(urls map[string]string, changeActions map[string]interface{}) error {
	functionName := "preparePRComment"

	// Generate pull request comment markdown
	commentMD, commentErr := ci.GenerateMD(urls, changeActions)
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
		"--showChanges", strconv.FormatBool(diagramValues["ShowChanges"].(bool)),
		"--showDrift", strconv.FormatBool(diagramValues["ShowDrift"].(bool)),
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

	// Capture change actions from Stdout
	var changeActions map[string]interface{}

	scanner := bufio.NewScanner(&outputSink)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "CHANGEACTIONS") {
			actionString := strings.ReplaceAll(line, "CHANGEACTIONS:", "")

			if unmarshalErr := json.Unmarshal([]byte(actionString), &changeActions); unmarshalErr != nil {
				return fmt.Errorf("parsing change actions failed -> %v: %w", functionName, unmarshalErr)
			}
		}
	}

	exportPath := filepath.Join(diagramValues["OutDir"].(string), diagramValues["FileName"].(string)) + ".pdf"

	// If environment is CI or --generate-md is set -> Host exported diagram and generate PR comment markdown
	if auxiliary.StateInstance.IsCI || diagramValues["GenerateMd"].(bool) {
		if auxiliary.StateInstance.PluralithConfig.ProjectId == "" {
			ux.PrintFormatted("\n✘", []string{"red", "bold"})
			fmt.Print(" No project ID set → Run ")
			ux.PrintFormatted("pluralith init", []string{"blue"})
			fmt.Println(" or provide a valid config\n")
			return nil
		}
		// Upload diagram to storage for pull request comment hosting
		urls, hostErr := PostRun(exportPath, changeActions)
		if hostErr != nil {
			return fmt.Errorf("hosting diagram for PR comment failed -> %v: %w", functionName, hostErr)
		}

		if prErr := GenerateComment(urls, changeActions); prErr != nil {
			return fmt.Errorf("handling pull request update failed -> %v: %w", functionName, prErr)
		}
	}
	// else {
	// 	if logErr := LogExport(); logErr != nil {
	// 		return fmt.Errorf("logging diagram export failed -> %v: %w", functionName, logErr)
	// 	}
	// }

	exportSpinner.Success()
	ux.PrintFormatted("  → ", []string{"blue"})
	fmt.Print("Diagram exported to: ")
	ux.PrintFormatted(exportPath, []string{"blue"})
	fmt.Println("\n")

	return nil
}
