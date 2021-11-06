package stream

import (
	"io/ioutil"
	"os"
	"path"
	"pluralith/pkg/auxiliary"
	"strings"
	"time"
)

func FetchResource(address string, isDestroy bool) (string, error) {
	// Define working dir
	workingDir, workingErr := os.Getwd()
	if workingErr != nil {
		return "", workingErr
	}
	// Define variable to hold state string
	var stateString string

	// Catching the terraform state right on an apply event is a bit all over the place
	// Sometimes the terraform.tfstate file is already updated, sometimes it isn't
	// Therefore we use a while loop to continuously poll terraform.tfstate and check for the presence of the recently updated resource's name
	// On apply -> should be present (should not equal isDestroy), on destroy -> shouldn't be present anymore (should equal isDestroy)
	for strings.Contains(stateString, strings.Split(address, ".")[1]) == isDestroy {
		// Introduce delay to avoid unnecessarily aggressive polling
		time.Sleep(10 * time.Millisecond)
		// Read tfstate file
		stateBytes, stateErr := ioutil.ReadFile(path.Join(workingDir, "terraform.tfstate"))
		// Convert read state to string
		stateString = string(stateBytes) // Assign
		if stateErr != nil {
			return "", stateErr
		}
	}

	// Parsing state object
	_, parseErr := auxiliary.ParseJson(stateString)
	if parseErr != nil {
		return "", parseErr
	}

	return "", nil
}
