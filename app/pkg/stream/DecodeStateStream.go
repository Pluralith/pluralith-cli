package stream

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"strings"
)

func DecodeStateStream(jsonString string) (string, string, error) {
	functionName := "DecodeStateStream"

	// Parsing state stream JSON
	parsedState, parseErr := auxiliary.ParseJson(jsonString)
	if parseErr != nil {
		return "", "", fmt.Errorf("could not parse json -> %v: %w", functionName, parseErr)
	}

	// Retrieving event type from parsed state JSON
	eventType := parsedState["type"].(string)

	// Filtering for "apply" event types
	if strings.Contains(eventType, "apply") {
		// Getting address of current resource
		hook := parsedState["hook"].(map[string]interface{})
		resource := hook["resource"].(map[string]interface{})
		address := resource["addr"]

		// If address key is given in current object -> Return event type and address
		if address != nil {
			return eventType, address.(string), nil
		}
	}

	return eventType, "", nil
}
