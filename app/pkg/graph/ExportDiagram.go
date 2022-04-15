package graph

import (
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

		// Read cache from disk
		cacheByte, cacheErr := os.ReadFile(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.cache.json"))
		if cacheErr != nil {
			return fmt.Errorf("reading cache from disk failed -> %v: %w", functionName, cacheErr)
		}

		// Unmarshal cache
		var runCache map[string]interface{}
		if unmarshallErr := json.Unmarshal(cacheByte, &runCache); unmarshallErr != nil {
			return fmt.Errorf("unmarshalling cache failed -> %v: %w", functionName, unmarshallErr)
		}

		// Upload diagram to storage for pull request comment hosting
		runData, postErr := PostRun(exportPath)
		if postErr != nil {
			return fmt.Errorf("posting run for PR comment failed -> %v: %w", functionName, postErr)
		}

		// Populate run cache data with additional attributes
		runCache["id"] = runData["id"]
		runCache["urls"] = runData["urls"]
		runCache["source"] = "CI"

		logErr := LogRun(runCache)
		if logErr != nil {
			return fmt.Errorf("posting run for PR comment failed -> %v: %w", functionName, logErr)
		}

		if prErr := GenerateComment(runCache); prErr != nil {
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
