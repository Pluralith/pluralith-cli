package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
)

func SetAPIKey(APIKey string) error {
	functionName := "SetAPIKey"

	auxiliary.StateInstance.APIKey = APIKey
	credentialsPath := filepath.Join(auxiliary.StateInstance.PluralithPath, "credentials")

	// Write api key to credentials file
	if writeErr := os.WriteFile(credentialsPath, []byte(APIKey), 0700); writeErr != nil {
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	return nil
}
