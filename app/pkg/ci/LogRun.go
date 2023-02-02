package ci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pluralith/pkg/auxiliary"
)

func LogRun(runCache map[string]interface{}, exportArgs map[string]interface{}) error {
	functionName := "LogRun"

	// Stringify run cache
	runCacheBytes, marshalErr := json.MarshalIndent(runCache, "", " ")
	if marshalErr != nil {
		return fmt.Errorf("creating run cache string failed -> %v: %w", functionName, marshalErr)
	}

	request, _ := http.NewRequest("POST", auxiliary.StateInstance.PluralithConfig.PluralithAPIEndpoint+"/v1/run/log", bytes.NewBuffer(runCacheBytes))
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return fmt.Errorf("logging run details failed -> %v: %w", functionName, responseErr)
	}

	// Parse response for file URLs
	responseBody, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return fmt.Errorf("%v: %w", functionName, readErr)
	}

	var bodyObject map[string]interface{}
	parseErr := json.Unmarshal(responseBody, &bodyObject)
	if parseErr != nil {
		return fmt.Errorf("parsing response failed -> %v: %w", functionName, parseErr)
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("request failed -> %v: %v", functionName, bodyObject)
	}

	runCache["urls"] = bodyObject["data"].(map[string]interface{})["urls"]
	exportArgs["run-id"] = bodyObject["data"].(map[string]interface{})["runId"].(string)

	return nil
}
