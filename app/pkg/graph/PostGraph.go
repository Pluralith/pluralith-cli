package graph

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func PostGraph(runType string, exportArgs map[string]interface{}) error {
	functionName := "PostGraph"

	fmt.Println()
	runSpinner := ux.NewSpinner("Posting Diagram", "Diagram Posted", "Posting Diagram Failed", true)
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
	runCache["id"] = exportArgs["run-id"]
	runCache["projectId"] = exportArgs["project-id"]
	runCache["type"] = "plan"
	runCache["version"] = ""
	runCache["branch"] = "local"

	config := make(map[string]interface{})
	config["showChanges"] = exportArgs["show-changes"]
	config["showCosts"] = exportArgs["show-costs"]
	config["showDrift"] = exportArgs["show-drift"]
	config["format"] = "png"

	runCache["config"] = config

	logErr := LogLocalRun(runCache, exportArgs)
	if logErr != nil {
		runSpinner.Fail()
		return fmt.Errorf("posting diagram to pluralith Dashboard failed -> %v: %w", functionName, logErr)
	}

	runSpinner.Success()
	ux.PrintFormatted("  → ", []string{"blue"})
	fmt.Print("Diagram Pushed To: ")
	ux.PrintFormatted(runCache["urls"].(map[string]interface{})["pluralithURL"].(string)+"\n", []string{"blue"})

	// Open run in deafult browser
	auxiliary.OpenBrowser(runCache["urls"].(map[string]interface{})["pluralithURL"].(string))

	return nil
}
