package stream

import (
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/comdb"
	"strings"
	"time"
)

func ProcessTerraformMessage(message string, command string) DecodedEvent {
	// functionName := "DecodeStateStream"
	decodedEvent := DecodedEvent{}

	// Parsing terraform message
	parsedState, parseErr := auxiliary.ParseJson(message)
	if parseErr != nil {
		return decodedEvent // If message is not valid json that cannot be parsed -> return and do nothing
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

		address := resource["addr"].(string)

		if address != "" {
			address = strings.Replace(address, "module.", "", 1)
			address = strings.Replace(address, "[", ".", 1)
			address = strings.Replace(address, "]", "", 1)
		}

		// Set address and type
		decodedEvent.Command = command
		decodedEvent.Address = address
		decodedEvent.Type = strings.Split(eventType, "_")[1]
	}

	// Handle diagnostic events
	if eventType == "diagnostic" {
		// Get address of current resource
		diagnostic := parsedState["diagnostic"].(map[string]interface{})
		eventType := parsedState["@level"].(string)

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

	// If address is given -> Resource event
	if decodedEvent.Address != "" {
		// Emit current event update to UI
		comdb.PushComDBEvent(comdb.ComDBEvent{
			Receiver:  "UI",
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Command:   decodedEvent.Command,
			Type:      decodedEvent.Type,
			Address:   decodedEvent.Address,
			Message:   decodedEvent.Message,
			Path:      auxiliary.StateInstance.WorkingPath,
			Received:  false,
		})
	}

	return decodedEvent
}
