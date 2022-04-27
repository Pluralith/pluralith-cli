package auth

import (
	"fmt"
	"net/http"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyProject(projectId string) (bool, error) {
	functionName := "VerifyAPIKey"

	verificationSpinner := ux.NewSpinner("Verifying Project ID", "Project ID is valid!", "No project with this ID exists, try again!", true)
	verificationSpinner.Start()

	// Construct key verification request
	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/project/get", nil)
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)

	// Add project id query string
	queryString := request.URL.Query()
	queryString.Add("projectId", projectId)
	request.URL.RawQuery = queryString.Encode()

	// Execute key verification request
	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Handle verification response
	if response.StatusCode == 200 {
		verificationSpinner.Success()
		return true, nil
	} else if response.StatusCode == 401 {
		verificationSpinner.Fail("You are not authorized to access this project")
		return false, nil
	} else {
		verificationSpinner.Fail()
		return false, nil
	}
}
