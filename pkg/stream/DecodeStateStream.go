package stream

import (
	"errors"
	"pluralith/pkg/auxiliary"
	"strings"
)

func DecodeStateStream(jsonString string) (string, string, error) {
	// Parsing state stream JSON
	parsedState, parseErr := auxiliary.ParseJson(jsonString)
	if parseErr != nil {
		return "", "", errors.New("decode_state_stream: Could not parse json")
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
