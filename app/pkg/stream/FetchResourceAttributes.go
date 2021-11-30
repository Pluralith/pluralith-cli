package stream

import (
	"strings"
)

func FetchResourceInstances(address string, stateObject map[string]interface{}) ([]interface{}, error) {
	// Find match based on resource name
	for _, item := range stateObject["resources"].([]interface{}) {
		name := strings.Split(address, ".")[1]
		itemMap := item.(map[string]interface{})

		// If resource with name present
		if itemMap["name"] == name {
			return itemMap["instances"].([]interface{}), nil
		}
	}

	return make([]interface{}, 0), nil
}
