package graph

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func HandleCIRun(exportArgs map[string]interface{}) error {
	functionName := "HandleCIRun"

	// Check for missing project ID
	if auxiliary.StateInstance.PluralithConfig.ProjectId == "" {
		ux.PrintFormatted("\n✘", []string{"red", "bold"})
		fmt.Print(" No project ID set → Run ")
		ux.PrintFormatted("pluralith init", []string{"blue"})
		fmt.Println(" or provide a valid config\n")
		return nil
	}

	runSpinner := ux.NewSpinner("Posting Run", "Run Posted To Pluralith Dashboard", "Posting Run To Pluralith Dashboard Failed", true)
	runSpinner.Start()

	// Read cache from disk
	cacheByte, cacheErr := os.ReadFile(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.cache.json"))
	if cacheErr != nil {
		runSpinner.Fail()
		return fmt.Errorf("reading cache from disk failed -> %v: %w", functionName, cacheErr)
	}

	// Unmarshal cache
	var runCache map[string]interface{}
	if unmarshallErr := json.Unmarshal(cacheByte, &runCache); unmarshallErr != nil {
		runSpinner.Fail()
		return fmt.Errorf("unmarshalling cache failed -> %v: %w", functionName, unmarshallErr)
	}

	// Populate run cache data with additional attributes
	runCache["id"] = exportArgs["id"]
	runCache["branch"] = exportArgs["branch"]

	config := make(map[string]interface{})
	config["showChanges"] = exportArgs["show-changes"]
	config["showCosts"] = exportArgs["show-costs"]
	config["showDrift"] = exportArgs["show-drift"]
	config["format"] = "png"

	runCache["config"] = config

	logErr := LogRun(runCache)
	if logErr != nil {
		runSpinner.Fail()
		return fmt.Errorf("posting run for PR comment failed -> %v: %w", functionName, logErr)
	}

	if prErr := GenerateComment(runCache); prErr != nil {
		runSpinner.Fail()
		return fmt.Errorf("handling pull request update failed -> %v: %w", functionName, prErr)
	}

	runSpinner.Success()

	return nil
}
