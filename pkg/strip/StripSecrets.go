package strip

import auxiliary "pluralith/pkg/auxiliary"

// - - - Code to strip secrets from provided JSON input - - -

// Function to strip state of secrets
func StripSecrets(jsonString string, targets []string, replacement string) (string, error) {
	// // Initializing empty variable to unmarshal JSON into
	// var jsonObject map[string]interface{}
	// // Unmarshalling JSON and handling potential errors
	// if err := json.Unmarshal([]byte(jsonString), &jsonObject); err != nil {
	// 	fmt.Println("Failure String\n", jsonString)
	// 	return "", errors.New("strip json unmarshal failed")
	// }

	processedJson, processErr := auxiliary.ProcessJson(jsonString)
	if processErr != nil {
		return "", processErr
	}

	// Calling recursive function to strip secrets and replace values on every level in JSON
	ReplaceSensitive(processedJson, targets, replacement)
	// Properly formating returned JSON
	// strippedObject, err := json.MarshalIndent(jsonObject, "", " ")
	// if err != nil {
	// 	return "", errors.New("stripping secrets failed")
	// } else {
	// 	return string(strippedObject), nil
	// }
}
