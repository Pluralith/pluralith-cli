package initialization

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyProject(orgId string, projectId string) (bool, string, error) {
	functionName := "VerifyProject"
	verificationResponse := ProjectResponse{}

	verificationSpinner := ux.NewSpinner("Verifying Project ID", "Project ID Is Valid", "No project with this ID exists, try again!", true)
	verificationSpinner.Start()

	if projectId == "" {
		verificationSpinner.Fail("No Project ID given")
		return false, "", nil
	}

	// Construct key verification request
	// request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/project/get", nil)
	request, _ := http.NewRequest("GET", "http://localhost:8080/v1/project/get", nil)
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
		return false, "", fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Parse response for file URLs
	responseBody, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return false, "", fmt.Errorf("%v: %w", functionName, readErr)
	}

	parseErr := json.Unmarshal(responseBody, &verificationResponse)
	if parseErr != nil {
		return false, "", fmt.Errorf("parsing response failed -> %v: %w", functionName, parseErr)
	}

	if response.StatusCode == 200 {
		verificationSpinner.Success("Existing Project Found")
		return true, verificationResponse.Data.Name, nil
	}

	if response.StatusCode == 404 {
		verificationSpinner.Success("No Project Found")
		return true, "", nil
	}

	if response.StatusCode == 401 {
		verificationSpinner.Fail("Not Authorized To Access This Project")
		return false, "", nil
	}

	verificationSpinner.Fail()
	return false, "", nil
}
