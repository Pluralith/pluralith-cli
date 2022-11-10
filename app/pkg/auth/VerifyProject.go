package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyProject(orgId string, projectId string) (ProjectResponse, error) {
	functionName := "VerifyProject"
	verificationResponse := ProjectResponse{}

	verificationSpinner := ux.NewSpinner("Verifying Project ID", "Project ID is valid!", "No project with this ID exists, try again!", true)
	verificationSpinner.Start()

	// Construct key verification request
	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/project/get", nil)
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)

	// Add project id query string
	queryString := request.URL.Query()
	queryString.Add("orgId", orgId)
	queryString.Add("projectId", projectId)
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
		verificationSpinner.Success("Existing Project Found")
		return verificationResponse, nil
	} else if response.StatusCode == 404 {
		verificationSpinner.Fail("No Project Found")
		return verificationResponse, nil
	} else if response.StatusCode == 401 {
		verificationSpinner.Fail("Not Authorized To Access This Project")
		return verificationResponse, nil
	}

	verificationSpinner.Fail()
	return verificationResponse, nil
}
