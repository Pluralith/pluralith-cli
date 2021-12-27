package stream

import (
	"fmt"
	"pluralith/pkg/auxiliary"
	"strings"
)

func DecodeStateStream(jsonString string) (DecodedEvent, error) {
	functionName := "DecodeStateStream"
	decodedEvent := DecodedEvent{}

	// Parsing state stream JSON
	parsedState, parseErr := auxiliary.ParseJson(jsonString)
	if parseErr != nil {
		return decodedEvent, fmt.Errorf("could not parse json -> %v: %w", functionName, parseErr)
	}

	// Retrieving event type from parsed state JSON
	decodedEvent.Type = parsedState["type"].(string)
	decodedEvent.Message = parsedState["@message"].(string)

	// Filtering for "apply" event types
	if strings.Contains(decodedEvent.Type, "apply") {
		// Getting address of current resource
		hook := parsedState["hook"].(map[string]interface{})
		resource := hook["resource"].(map[string]interface{})
		decodedEvent.Address = resource["addr"].(string)
	}

	// if decodedEvent

	return decodedEvent, nil
}
