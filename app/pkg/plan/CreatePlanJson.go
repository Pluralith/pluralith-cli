package plan

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"pluralith/pkg/auxiliary"
	"pluralith/pkg/strip"
	"pluralith/pkg/ux"
)

func SplitJsonPlan(planJsonString string) ([]string, error) {
	functionName := "SplitJsonPlan"

	var plans []string
	bracketsOpened := 0
	bracketsClosed := 0
	subString := ""
	objectStarted := false

	for _, char := range planJsonString {
		if char == '{' {
			objectStarted = true
			bracketsOpened++
		} else if char == '}' {
			bracketsClosed++
		}

		if objectStarted {
			subString += string(char)
		}

		if objectStarted && bracketsOpened > 0 && bracketsOpened == bracketsClosed {
			_, parseErr := auxiliary.ParseJson(subString)
			if parseErr != nil {
				return nil, fmt.Errorf("%v: %w", functionName, parseErr)
			}
			plans = append(plans, subString)
			subString = ""
			objectStarted = false
			bracketsClosed = 0
			bracketsOpened = 0
		}
	}

	return plans, nil
}

func ConvertBinaryPlanToJson(planPath string) (string, error) {
	functionName := "ConvertBinaryPlanToJson"

	// Constructing command to execute
	cmd := exec.Command("terraform", append([]string{"show", "-json", planPath})...)

	// Defining sinks for std data
	var outputSink bytes.Buffer
	var errorSink bytes.Buffer

	// Redirecting command std data
	cmd.Stdout = &outputSink
	cmd.Stderr = &errorSink
	cmd.Stdin = os.Stdin

	// Run terraform command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("terraform command failed -> %v: %w", functionName, err)
	}

	return outputSink.String(), nil
}

func CreatePlanJson(planPath string, isJson bool) (string, []string, []string, error) {
	functionName := "CreatePlanJson"

	// Initiate variables
	var jsonPlans []string
	var strippedPlans []string
	var providers []string

	// Construct file path for stripped state
	strippedPath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.state.json")

	// Instantiate spinner
	stripSpinner := ux.NewSpinner("Stripping Secrets", "Secrets Stripped", "Stripping Secrets Failed", true)
	stripSpinner.Start()

	planJsonString := ""
	if isJson {
		// Read json file
		file, readErr := ioutil.ReadFile(planPath)
		if readErr != nil {
			stripSpinner.Fail()
			return strippedPath, jsonPlans, providers, fmt.Errorf("failed to open json plan -> %v: %w", functionName, readErr)
		}
		planJsonString = string(file)
	} else {
		// Convert binary plan to json
		planJsonStringTemp, convertErr := ConvertBinaryPlanToJson(planPath)
		if convertErr != nil {
			stripSpinner.Fail()
			return strippedPath, jsonPlans, providers, fmt.Errorf("failed to strip secrets -> %v: %w", functionName, convertErr)
		}
		planJsonString = planJsonStringTemp
	}

	// Split json plans
	jsonPlans, splitError := SplitJsonPlan(planJsonString)
	if splitError != nil {
		stripSpinner.Fail()
		return strippedPath, jsonPlans, providers, fmt.Errorf("failed to split json plan -> %v: %w", functionName, splitError)
	}

	// Strip secrets from plan states json
	for _, jsonPlan := range jsonPlans {
		strippedPlan, stripErr := strip.StripSecrets(jsonPlan)
		if stripErr != nil {
			stripSpinner.Fail()
			return strippedPath, jsonPlans, providers, fmt.Errorf("failed to strip secrets -> %v: %w", functionName, stripErr)
		}
		strippedPlans = append(strippedPlans, strippedPlan)
	}

	// Fetch providers present in states json
	for _, jsonPlan := range jsonPlans {
		providersTemp, providerErr := FetchProviders(jsonPlan)
		if providerErr != nil {
			stripSpinner.Fail()
			return strippedPath, jsonPlans, providers, fmt.Errorf("failed to fetch providers -> %v: %w", functionName, providerErr)
		}
		providers = append(providers, providersTemp...)
	}

	// Merge plans
	strippedPlanString := "["
	for index, jsonPlan := range strippedPlans {
		strippedPlanString += jsonPlan
		if index < len(strippedPlans)-1 {
			strippedPlanString += ","
		}
	}
	strippedPlanString += "]"

	// Write stripped state to file
	if writeErr := ioutil.WriteFile(strippedPath, []byte(strippedPlanString), 0700); writeErr != nil {
		stripSpinner.Fail()
		return strippedPath, jsonPlans, providers, fmt.Errorf("%v: %w", functionName, writeErr)
	}

	stripSpinner.Success()

	// Return path to execution plan
	return strippedPath, jsonPlans, providers, nil
}
