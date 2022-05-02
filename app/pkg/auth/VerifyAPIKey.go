package auth

import (
	"fmt"
	"net/http"
	"pluralith/pkg/ux"
)

func VerifyAPIKey(APIKey string, silent bool) (bool, error) {
	functionName := "VerifyAPIKey"

	verificationSpinner := ux.NewSpinner("Verifying your API key", "API key is valid, you are authenticated!", "The passed API key is invalid, try again!", true)
	if !silent {
		verificationSpinner.Start()
	}

	// Construct key verification request
	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/auth/key/verify", nil)
	request.Header.Add("Authorization", "Bearer "+APIKey)

	// Execute key verification request
	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Handle verification response
	if response.StatusCode == 200 {
		if !silent {
			verificationSpinner.Success()
		}
		return true, nil
	} else {
		if !silent {
			verificationSpinner.Fail()
		}
		return false, nil
	}
}
