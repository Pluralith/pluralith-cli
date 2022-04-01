package cost

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
)

func CalculateCost() {
	costSpinner := ux.NewSpinner("Calculating infrastructure costs", "Costs calculated", "Couldn't calculate costs", true)
	costSpinner.Start()

	planJsonPath := filepath.Join(auxiliary.StateInstance.WorkingPath, "pluralith.state.stripped")

	costCmd := exec.Command("infracost", append([]string{"breakdown", "--path", planJsonPath, "--out-file", "cost.json", "--format", "json"})...)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	costCmd.Stdout = &outputSink
	costCmd.Stderr = &errorSink
	costCmd.Stdin = os.Stdin

	if runErr := costCmd.Run(); runErr != nil {
		fmt.Println(runErr)
		fmt.Println(errorSink.String())
		costSpinner.Fail()
	}

	costSpinner.Success()
}
