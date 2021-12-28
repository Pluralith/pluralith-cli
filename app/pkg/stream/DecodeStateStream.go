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

	// Get event message
	decodedEvent.Message = parsedState["@message"].(string)

	// Retrieve event type from parsed state JSON
	eventType := parsedState["type"].(string)

	// Handle apply events
	if strings.Contains(eventType, "apply") {
		// Get address of current resource
		hook := parsedState["hook"].(map[string]interface{})
		resource := hook["resource"].(map[string]interface{})

		// Set address and type
		decodedEvent.Address = resource["addr"].(string)
		decodedEvent.Type = strings.Split(eventType, "_")[1]
	}

	// Handle diagnostic events
	if eventType == "diagnostic" {
		// Get address of current resource
		diagnostic := parsedState["diagnostic"].(map[string]interface{})

		// Set address and type
		decodedEvent.Address = diagnostic["address"].(string)
		decodedEvent.Type = parsedState["@level"].(string)
	}

	return decodedEvent, nil
}
