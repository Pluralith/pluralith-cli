package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pluralith/pkg/auxiliary"
)

func LogRun(runCache map[string]interface{}) error {
	functionName := "LogRun"

	runCache["projectId"] = auxiliary.StateInstance.PluralithConfig.ProjectId

	// Stringify run cache
	runCacheBytes, marshalErr := json.MarshalIndent(runCache, "", " ")
	if marshalErr != nil {
		return fmt.Errorf("creating run cache string failed -> %v: %w", functionName, marshalErr)
	}

	// request, _ := http.NewRequest("POST", "https://api.pluralith.com/v1/run/log", bytes.NewBuffer(runCacheBytes))
	request, _ := http.NewRequest("POST", "http://localhost:8080/v1/run/log", bytes.NewBuffer(runCacheBytes))
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return fmt.Errorf("logging run details failed -> %v: %w", functionName, responseErr)
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("request failed -> %v: %v", functionName, response)
	}

	// Parse response for file URLs
	responseBody, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return fmt.Errorf("%v: %w", functionName, readErr)
	}

	var bodyObject map[string]interface{}
	parseErr := json.Unmarshal(responseBody, &bodyObject)
	if parseErr != nil {
		return fmt.Errorf("parsing response failed -> %v: %w", functionName, responseErr)
	}

	runCache["urls"] = bodyObject["data"]

	return nil
}
