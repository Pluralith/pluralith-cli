package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyOrg(orgId string) (OrgResponse, error) {
	functionName := "VerifyOrg"
	verificationResponse := OrgResponse{}

	verificationSpinner := ux.NewSpinner("Verifying Org ID", "Org ID is valid!", "No org with this ID exists, try again!", true)
	verificationSpinner.Start()

	// Construct key verification request
	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/org/get", nil)
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)

	// Add project id query string
	queryString := request.URL.Query()
	queryString.Add("orgId", orgId)
	request.URL.RawQuery = queryString.Encode()

	// Execute key verification request
	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		return verificationResponse, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Parse response for file URLs
	responseBody, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return verificationResponse, fmt.Errorf("%v: %w", functionName, readErr)
	}

	parseErr := json.Unmarshal(responseBody, &verificationResponse)
	if parseErr != nil {
		return verificationResponse, fmt.Errorf("parsing response failed -> %v: %w", functionName, parseErr)
	}

	if response.StatusCode == 200 {
		verificationSpinner.Success("Org Found")
		return verificationResponse, nil
	}

	verificationSpinner.Fail()
	return verificationResponse, nil
}

// func VerifyOrg(orgId string) (map[string]interface{}, error) {
// 	functionName := "VerifyProject"

// 	verificationSpinner := ux.NewSpinner("Verifying Project ID", "Project ID is valid!", "No project with this ID exists, try again!", true)
// 	verificationSpinner.Start()

// 	// Construct key verification request
// 	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/org/get", nil)
// 	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)

// 	// Add org id query string
// 	queryString := request.URL.Query()
// 	queryString.Add("orgId", orgId)
// 	request.URL.RawQuery = queryString.Encode()

// 	// Execute key verification request
// 	client := &http.Client{}
// 	response, responseErr := client.Do(request)
// 	if responseErr != nil {
// 		return nil, fmt.Errorf("%v: %w", functionName, responseErr)
// 	}

// 	// Parse response for file URLs
// 	responseBody, readErr := ioutil.ReadAll(response.Body)
// 	if readErr != nil {
// 		return nil, fmt.Errorf("%v: %w", functionName, readErr)
// 	}

// 	var bodyObject map[string]interface{}
// 	parseErr := json.Unmarshal(responseBody, &bodyObject)
// 	if parseErr != nil {
// 		return nil, fmt.Errorf("parsing response failed -> %v: %w", functionName, parseErr)
// 	}

// 	// Handle verification response
// 	if response.StatusCode == 200 {
// 		verificationSpinner.Success()
// 		return bodyObject, nil
// 	} else if response.StatusCode == 401 {
// 		verificationSpinner.Fail("You are not authorized to access this org")
// 		return nil, nil
// 	} else {
// 		verificationSpinner.Fail()
// 		return nil, nil
// 	}
// }
