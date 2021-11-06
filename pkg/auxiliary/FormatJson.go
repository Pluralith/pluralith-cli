package auxiliary

import (
	"encoding/json"
	"errors"
)

func FormatJson(jsonObject map[string]interface{}) (string, error) {
	// Properly format JSON string
	indented, indentErr := json.MarshalIndent(jsonObject, "", " ")
	if indentErr != nil {
		return "", errors.New("json formatting failed")
	}

	// Return result
	return string(indented), nil
}
