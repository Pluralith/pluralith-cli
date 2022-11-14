package auth

import "fmt"

func RunLogin(APIKey string) (bool, error) {
	functionName := "RunLogin"

	// Verify API key with backend
	isValid, verifyErr := VerifyAPIKey(APIKey, false)
	if verifyErr != nil {
		return false, fmt.Errorf("verifying API key failed -> %v: %w", functionName, verifyErr)
	}

	if isValid {
		// Set API key in credentials file at ~/Pluralith/credentials
		setErr := SetAPIKey(APIKey)
		if setErr != nil {
			return false, fmt.Errorf("setting API key in credentials file failed -> %v: %w", functionName, setErr)
		}
	} else {
		return false, nil
	}

	return true, nil
}
