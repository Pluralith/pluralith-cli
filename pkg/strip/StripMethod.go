package strip

import (
	"fmt"
	"io/ioutil"

	"pluralith/pkg/ux"
)

// Specifying sensitive keys (will later be done via external config)
var sensitiveKeys = []string{"tags", "owner_id"}

func StripMethod(args []string) {
	// Fetching all state files in current working directory
	stateFiles := FetchFiles(".tfstate")

	// Instantiating new strip spinner
	stripSpinner := ux.NewSpinner("Stripping Secrets", fmt.Sprintf("Secrets Stripped From %d File", len(stateFiles)), "Stripping Secrets Failed")
	stripSpinner.Start()

	// Stripping secrets and writing stripped state to disk
	for fileName, fileContent := range stateFiles {
		strippedFile, err := StripSecrets(fileContent, sensitiveKeys, "gatewatch")
		if err != nil {
			stripSpinner.Fail("Failed to strip secrets from %s", fileName)
		} else {
			ioutil.WriteFile(fmt.Sprintf("%s.state.stripped", fileName), []byte(strippedFile), 0644)
			stripSpinner.Success()
		}
	}
}
