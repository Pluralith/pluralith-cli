package graph

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func CheckValidity(version string, payload string) (bool, error) {
	functionName := "CheckValidity"

	// Compute checksum
	hasher := sha256.New()
	hasher.Write([]byte(payload))
	checksum := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	var bodyObject map[string]interface{}

	// Submit checksum for validation check
	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/graph/validate/"+version, nil)
	request.Header.Add("Authorization", "Bearer "+os.Getenv("PluralithDevApiKey"))

	queryString := request.URL.Query()
	queryString.Add("checksum", checksum)
	request.URL.RawQuery = queryString.Encode()

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Parse request body
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	if parseErr := json.Unmarshal(bodyBytes, &bodyObject); parseErr != nil {
		fmt.Println("unmarshal fails")
		return false, fmt.Errorf("%v: %w", functionName, parseErr)
	}

	return bodyObject["data"].(bool), nil
}
