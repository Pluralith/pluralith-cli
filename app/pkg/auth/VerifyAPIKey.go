package auth

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func VerifyAPIKey() (bool, error) {
	functionName := "VerifyAPIKey"

	verificationSpinner := ux.NewSpinner("Verifying your API key", "Your API key is valid, you are logged in!\n", "API key verification failed\n", false)

	ux.PrintFormatted("→", []string{"blue", "bold"})
	fmt.Print(" Enter API Key: ")

	// Capture user input
	var APIKey string
	fmt.Scanln(&APIKey)

	verificationSpinner.Start()

	// Construct key verification request
	request, _ := http.NewRequest("GET", "https://api.pluralith.com/v1/auth/key/verify", nil)
	request.Header.Add("Authorization", "Bearer "+APIKey)

	// Execute key verification request
	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		verificationSpinner.Fail("Failed to verify API key\n")
		return false, fmt.Errorf("%v: %w", functionName, responseErr)
	}

	// Hande verification response
	if response.StatusCode == 200 {
		auxiliary.StateInstance.APIKey = APIKey
		credentialsPath := filepath.Join(auxiliary.StateInstance.PluralithPath, "credentials")

		// Write api key to credentials file
		if writeErr := os.WriteFile(credentialsPath, []byte(APIKey), 0700); writeErr != nil {
			verificationSpinner.Fail("Failed to write API key to config try again!\n")
			return false, fmt.Errorf("%v: %w", functionName, writeErr)
		}

		verificationSpinner.Success()
	} else {
		verificationSpinner.Fail("The passed API key is invalid, try again!\n")
	}

	return true, nil
}
