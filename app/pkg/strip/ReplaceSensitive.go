package strip

import (
	"reflect"

	"pluralith/pkg/auxiliary"
)

// Function to recursively replace key values in JSON
func ReplaceSensitive(jsonObject map[string]interface{}) {
	// Iterate over current level key value pairs
	for outerKey, outerValue := range jsonObject {
		// Check if value at key is given
		if outerValue != nil {
			// Subsitute value with replacement if key is among targets
			if auxiliary.ElementInSlice(outerKey, auxiliary.StateInstance.PluralithConfig.Config.SensitiveAttrs) {
				jsonObject[outerKey] = "gatewatch"
			} else {
				// Get value type to handle different cases
				outerValueType := reflect.TypeOf(outerValue)

				// Switch between different data types
				switch outerValueType.Kind() {
				case reflect.Map:
					// If value is of type map -> Move on to next recursion level
					ReplaceSensitive(outerValue.(map[string]interface{}))
				case reflect.Array, reflect.Slice:
					// If value is of type array or slice -> Loop through elements, if maps are found -> Move to next recursion level
					for _, innerValue := range outerValue.([]interface{}) {
						if innerValue != nil && reflect.TypeOf(innerValue).Kind() == reflect.Map {
							ReplaceSensitive(innerValue.(map[string]interface{}))
						}
					}
				}
			}
		}
	}
}
