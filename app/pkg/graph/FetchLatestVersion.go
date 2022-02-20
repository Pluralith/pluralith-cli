package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func FetchLatestVersion() (string, error) {
	functionName := "fetchLatestVersion"
	// Initialize get version request
	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/graph/get/current/version", nil)
	request.Header.Add("Authorization", "Bearer "+os.Getenv("PluralithDevApiKey"))

	// Execute get version request
	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return "", fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Parse request body
	var bodyObject map[string]interface{}
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	parseErr := json.Unmarshal(bodyBytes, &bodyObject)
	if parseErr != nil {
		return "", fmt.Errorf("%v: %w", functionName, parseErr)
	}

	return bodyObject["data"].(string), nil
}
