package auxiliary

import (
	"encoding/json"
	"errors"
)

func ParseJson(jsonInput string) (map[string]interface{}, error) {
	// Define empty object to fill with unmarshalled content
	var parsed map[string]interface{}

	// Unmarshall JSON string
	parseErr := json.Unmarshal([]byte(jsonInput), &parsed)
	if parseErr != nil {
		return make(map[string]interface{}), errors.New("json unmarshal failed")
	}

	// Return result
	return parsed, nil
}
