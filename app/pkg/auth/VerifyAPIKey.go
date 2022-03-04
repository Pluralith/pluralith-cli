package auth

import (
	"fmt"
	"net/http"
)

func VerifyAPIKey(APIKey string) (bool, error) {
	functionName := "VerifyAPIKey"

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
		return true, nil
	} else {
		return false, nil
	}
}
