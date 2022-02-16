package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"pluralith/pkg/auxiliary"
)

func FetchLatestAlgorithm(cachePath string) (map[string]interface{}, error) {
	functionName := "fetchLatestAlgorithm"

	version := "0.0.3"

	var bodyObject map[string]interface{}

	request, _ := http.NewRequest("GET", "http://localhost:8080/v1/graph/get/"+version, nil)
	request.Header.Add("Authorization", "Bearer "+os.Getenv("PluralithDevApiKey"))

	queryString := request.URL.Query()
	queryString.Add("apiKey", auxiliary.StateInstance.APIKey)
	request.URL.RawQuery = queryString.Encode()

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

	// Decrypt payload
	decryptedPayload, decryptErr := DecryptPayload(bodyObject["payload"].(string), bodyObject["iv"].(string))
	if decryptErr != nil {
		fmt.Println("crypto to the moon")
		return bodyObject, fmt.Errorf("failed to decrypt payload -> %v: %w", functionName, decryptErr)
	}

	// Check validity of payload
	isValid, checkErr := CheckValidity(version, decryptedPayload)
	if checkErr != nil {
		return bodyObject, fmt.Errorf("failed to check graphing validity -> %v: %w", functionName, checkErr)
	}

	if isValid {
		bodyObject["algorithm"] = decryptedPayload

		// Format cache string
		marshalString, marshalErr := json.MarshalIndent(bodyObject, "", " ")
		if marshalErr != nil {
			return bodyObject, fmt.Errorf("%v: %w", functionName, marshalErr)
		}

		// Write to local cache
		if writeErr := os.WriteFile(cachePath, marshalString, 0700); writeErr != nil {
			return bodyObject, fmt.Errorf("%v: %w", functionName, writeErr)
		}

		fmt.Println("Graphing updated to version " + bodyObject["version"].(string))
		return bodyObject, nil
	}

	fmt.Println("Received graphing data not valid")
	return bodyObject, nil
}
