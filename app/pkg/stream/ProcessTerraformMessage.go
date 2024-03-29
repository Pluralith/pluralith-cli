package stream

import (
	"pluralith/pkg/auxiliary"
	"strings"
)

func ProcessTerraformMessage(message string, command string) DecodedEvent {
	// functionName := "DecodeStateStream"
	decodedEvent := DecodedEvent{}

	// Parsing terraform message
	parsedMessage, parseErr := auxiliary.ParseJson(message)
	if parseErr != nil {
		return decodedEvent // If message is not valid json that cannot be parsed -> return and do nothing
	}

	// Get event message
	decodedEvent.Message = parsedMessage["@message"].(string)
	decodedEvent.Type = parsedMessage["type"].(string)

	// Handle apply events
	if strings.Contains(decodedEvent.Type, "apply") {
		// Get address of current resource
		hook := parsedMessage["hook"].(map[string]interface{})
		resource := hook["resource"].(map[string]interface{})

		address := resource["addr"].(string)

		if address != "" {
			address = strings.Replace(address, "module.", "", 1)
			address = strings.Replace(address, "[", ".", 1)
			address = strings.Replace(address, "]", "", 1)
		}

		// Set address and type
		decodedEvent.Command = command
		decodedEvent.Address = address
		decodedEvent.ParsedType = strings.Split(decodedEvent.Type, "_")[1]
	}

	// Handle diagnostic events
	if decodedEvent.Type == "diagnostic" {
		// Get address of current resource
		diagnostic := parsedMessage["diagnostic"].(map[string]interface{})
		eventType := parsedMessage["@level"].(string)

		if eventType == "error" {
			eventType = "errored"
		}

		// Set address and type
		if diagnostic["address"] != nil {
			decodedEvent.Command = "diagnostic"
			decodedEvent.Address = diagnostic["address"].(string)
			decodedEvent.Type = eventType
		}
	}

	// Handle output events
	if decodedEvent.Type == "outputs" {
		decodedEvent.Outputs = parsedMessage["outputs"].(map[string]interface{})
	}

	return decodedEvent
}
