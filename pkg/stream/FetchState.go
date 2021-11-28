package stream

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"pluralith/pkg/auxiliary"
)

func FetchState(address string, isDestroy bool) (map[string]interface{}, bool, error) {
	// Define working dir
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return make(map[string]interface{}), false, workingErr
	}
	// Define variables
	var stateString string
	var resourceFound bool = false
	var retries int = 0

	// Catching the terraform state right on an apply event is a bit all over the place
	// Sometimes the terraform.tfstate file is already updated, sometimes it isn't
	// Therefore we use a while loop to continuously poll terraform.tfstate and check for the presence of the recently updated resource's name
	// On apply -> should be present (should not equal isDestroy), on destroy -> shouldn't be present anymore (should equal isDestroy)
	for resourceFound == isDestroy && retries <= 20 {
		// Check if resource name has been found in string
		resourceFound = strings.Contains(stateString, strings.Split(address, ".")[1])
		// Introduce delay to avoid unnecessarily aggressive polling
		time.Sleep(20 * time.Millisecond)
		// Read tfstate file
		stateBytes, stateErr := ioutil.ReadFile(path.Join(workingDir, "terraform.tfstate"))
		// Convert read state to string
		stateString = string(stateBytes) // Assign
		if stateErr != nil {
			return make(map[string]interface{}), false, stateErr
		}

		retries += 1
	}

	// Parsing state object
	parsedState, parseErr := auxiliary.ParseJson(stateString) // FAILS ON DESTROY -> FIX
	if parseErr != nil {
		return make(map[string]interface{}), false, parseErr
	}

	return parsedState, resourceFound, nil
}
