package strip

import (
	"encoding/json"
	"fmt"
	"pluralith/pkg/auxiliary"
)

// Function to strip state of secrets
func StripSecrets(jsonString string, targets []string, replacement string) (string, error) {
	functionName := "StripSecrets"
	// Parse JSON object from string
	parsedJson, parseErr := auxiliary.ParseJson(jsonString)
	if parseErr != nil {
		return "", fmt.Errorf("%v: %w", functionName, parseErr)
	}

	// Call recursive function to strip secrets and replace values on every level in JSON
	ReplaceSensitive(parsedJson, targets, replacement)

	// Format JSON for output
	formattedJson, formatErr := json.MarshalIndent(parsedJson, "", " ")
	if formatErr != nil {
		return "", fmt.Errorf("%v: %w", functionName, formatErr)
	}

	// Return result
	return string(formattedJson), nil
}
