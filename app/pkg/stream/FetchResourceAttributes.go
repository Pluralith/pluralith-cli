package stream

import (
	"pluralith/pkg/strip"
	"strings"
)

func FetchResourceInstances(address string, stateObject map[string]interface{}) []interface{} {
	// Find match based on resource name
	for _, item := range stateObject["resources"].([]interface{}) {
		name := strings.Split(address, ".")[1]
		itemMap := item.(map[string]interface{})

		// Filter secrets according to config
		strip.ReplaceSensitive(itemMap)

		// If resource with name present
		if itemMap["name"] == name {
			return itemMap["instances"].([]interface{})
		}
	}

	return make([]interface{}, 0)
}
