package strip

import (
	"reflect"

	auxiliary "pluralith/pkg/auxiliary"
)

// Function to recursively replace key values in JSON
func ReplaceSensitive(jsonObject map[string]interface{}, targets []string, replacement string) {
	// Iterating over current level key value pairs
	for outerKey, outerValue := range jsonObject {
		// Checking if value at key is given
		if outerValue != nil {
			// Subsituting value with replacement if key is among targets
			if auxiliary.ElementInSlice(outerKey, targets) {
				jsonObject[outerKey] = replacement
			} else {
				// Getting value type to handle different cases
				outerValueType := reflect.TypeOf(outerValue)

				// Switching between different data types
				switch outerValueType.Kind() {
				case reflect.Map:
					// If value is of type map -> Move on to next recursion level
					ReplaceSensitive(outerValue.(map[string]interface{}), targets, replacement)
				case reflect.Array, reflect.Slice:
					// If value is of type array or slice -> Loop through elements, if maps are found -> Move to next recursion level
					for _, innerValue := range outerValue.([]interface{}) {
						if reflect.TypeOf(innerValue).Kind() == reflect.Map {
							ReplaceSensitive(innerValue.(map[string]interface{}), targets, replacement)
						}
					}
				}
			}
		}
	}
}
