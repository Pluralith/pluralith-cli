package strip

import auxiliary "pluralith/pkg/auxiliary"

// Function to strip state of secrets
func StripSecrets(jsonString string, targets []string, replacement string) (string, error) {
	// Parse JSON object from string
	parsedJson, parseErr := auxiliary.ParseJson(jsonString)
	if parseErr != nil {
		return "", parseErr
	}

	// Call recursive function to strip secrets and replace values on every level in JSON
	ReplaceSensitive(parsedJson, targets, replacement)

	// Format JSON for output
	formattedJson, formatErr := auxiliary.FormatJson(parsedJson)
	if formatErr != nil {
		return "", formatErr
	}

	// Return result
	return string(formattedJson), nil
}
