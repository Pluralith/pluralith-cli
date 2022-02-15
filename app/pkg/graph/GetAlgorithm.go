package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
)

func getInstalledVersion(cachePath string) (map[string]interface{}, error) {
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

func fetchLatestVersion() (string, error) {
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

func fetchLatestAlgorithm(cachePath string) (map[string]interface{}, error) {
	functionName := "fetchLatestAlgorithm"

	var bodyObject map[string]interface{}

	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/graph/get/current", nil)
	request.Header.Add("Authorization", "Bearer "+os.Getenv("PluralithDevApiKey"))

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return bodyObject, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Parse request body
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	if parseErr := json.Unmarshal(bodyBytes, &bodyObject); parseErr != nil {
		return bodyObject, fmt.Errorf("%v: %w", functionName, parseErr)
	}

	// Extract algorithm data
	bodyObject = bodyObject["algorithm"].(map[string]interface{})

	// Format cache string
	marshalString, marshalErr := json.MarshalIndent(bodyObject, "", " ")
	if marshalErr != nil {
		return bodyObject, fmt.Errorf("%v: %w", functionName, marshalErr)
	}

	// Write to local cache
	if writeErr := os.WriteFile(cachePath, marshalString, 0700); writeErr != nil {
		return bodyObject, fmt.Errorf("%v: %w", functionName, writeErr)
	}

	fmt.Println("Algorithm updated to version " + bodyObject["version"].(string))

	return bodyObject, nil
}

func GetAlgorithm() (map[string]interface{}, error) {
	functionName := "GetAlgorithm"

	// Define path to graphing cache
	cachePath := filepath.Join(auxiliary.StateInstance.PluralithPath, "pluralithCache.json")
	// Get installed graphing version
	installedVersion, versionErr := getInstalledVersion(cachePath)
	if versionErr != nil {
		// If no version installed -> fetch latest
		latestAlgorithm, fetchErr := fetchLatestAlgorithm(cachePath)
		if fetchErr != nil {
			return latestAlgorithm, fmt.Errorf("failed to fetch latest graphing data -> %v: %w", functionName, fetchErr)
		}

		return latestAlgorithm, nil
	}

	// Get latest production version string
	latestVersion, versionErr := fetchLatestVersion()
	if versionErr != nil {
		fmt.Println(fmt.Errorf("failed to fetch latest graphing version -> %v: %w", functionName, versionErr))
	}

	// Compare installed and latest version
	if installedVersion["version"] != latestVersion {
		// If latest version not equal to installed version -> fetch latest
		latestAlgorithm, fetchErr := fetchLatestAlgorithm(cachePath)
		if fetchErr != nil {
			return latestAlgorithm, fmt.Errorf("%v: %w", functionName, fetchErr)
		}

		return latestAlgorithm, nil
	}

	fmt.Println("Algorithm version " + installedVersion["version"].(string) + " up to date")

	return installedVersion, nil
}
