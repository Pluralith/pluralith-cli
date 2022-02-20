package graph

import (
	"encoding/json"
	"fmt"
	"os"
)

func FetchInstalledAlgorithm(cachePath string) (map[string]interface{}, error) {
	functionName := "getCurrentVerion"

	var cacheObject map[string]interface{}

	// Read cache from filesystem
	cacheValue, readErr := os.ReadFile(cachePath)
	if readErr != nil {
		return cacheObject, fmt.Errorf("%v: %w", functionName, readErr)
	}

	// Unmarshall cache JSON object
	parseErr := json.Unmarshal(cacheValue, &cacheObject)
	if parseErr != nil {
		return cacheObject, fmt.Errorf("%v: %w", functionName, parseErr)
	}

	return cacheObject, nil
}
