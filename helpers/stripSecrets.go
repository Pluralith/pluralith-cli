package helpers

import (
	"encoding/json"
	"log"
	"reflect"
)

// - - - Code to strip secrets from provided JSON input - - -

// Function to recursively replace key values in JSON
func replaceSensitive(jsonObject map[string]interface{}, targets []string, replacement string) {
	// Iterating over current level key value pairs
	for outerKey, outerValue := range jsonObject {
		// Checking if value at key is given
		if outerValue != nil {
			// Subsituting value with replacement if key is among targets
			if ElementInSlice(outerKey, targets) {
				jsonObject[outerKey] = replacement
			} else {
				// Getting value type to handle different cases
				outerValueType := reflect.TypeOf(outerValue)

				// Switching between different data types
				switch outerValueType.Kind() {
				case reflect.Map:
					// If value is of type map -> Move on to next recursion level
					replaceSensitive(outerValue.(map[string]interface{}), targets, replacement)
				case reflect.Array, reflect.Slice:
					// If value is of type array or slice -> Loop through elements, if maps are found -> Move to next recursion level
					for _, innerValue := range outerValue.([]interface{}) {
						if reflect.TypeOf(innerValue).Kind() == reflect.Map {
							replaceSensitive(innerValue.(map[string]interface{}), targets, replacement)
						}
					}
				}
			}
		}
	}
}

// Function to strip state of secrets
func StripSecrets(jsonString string, targets []string, replacement string) string {
	// Initializing empty variable to unmarshal JSON into
	var jsonObject map[string]interface{}
	// Unmarshalling JSON and handling potential errors
	if err := json.Unmarshal([]byte(jsonString), &jsonObject); err != nil {
		log.Fatal("Unmarshalling JSON failed")
	}
	// Calling recursive function to strip secrets and replace values on every level in JSON
	replaceSensitive(jsonObject, targets, replacement)
	// Properly formating returned JSON
	cleanObject, err := json.MarshalIndent(jsonObject, "", " ")
	if err != nil {
		return "Stripping secrets failed"
	} else {
		return string(cleanObject)
	}
}
