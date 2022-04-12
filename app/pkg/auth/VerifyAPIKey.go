package auth

import (
	"fmt"
	"net/http"
	"pluralith/pkg/ux"
)

func VerifyAPIKey(APIKey string) (bool, error) {
	functionName := "VerifyAPIKey"

	verificationSpinner := ux.NewSpinner("Verifying your API key", "Your API key is valid, you are authenticated!\n", "The passed API key is invalid, try again!\n", true)
	verificationSpinner.Start()

	// Construct key verification request
	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/auth/key/verify", nil)
	request.Header.Add("Authorization", "Bearer "+APIKey)

	// Execute key verification request
	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return false, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Hande verification response
	if response.StatusCode == 200 {
		verificationSpinner.Success()
		return true, nil
	} else {
		verificationSpinner.Fail()
		return false, nil
	}
}
