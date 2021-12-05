package auxiliary

import (
	"encoding/json"
	"fmt"
)

func ParseJson(jsonInput string) (map[string]interface{}, error) {
	functionName := "ParseJson"
	// Define empty object to fill with unmarshalled content
	var parsed map[string]interface{}

	// Unmarshall JSON string
	parseErr := json.Unmarshal([]byte(jsonInput), &parsed)
	if parseErr != nil {
		return make(map[string]interface{}), fmt.Errorf("%v: %w", functionName, parseErr)
	}

	// Return result
	return parsed, nil
}
