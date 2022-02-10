package strip

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
	"reflect"
	"strconv"
	"strings"
)

// TODO:
// DONE - Remove tags, tags_all, description, source
// - Handle sensitive stuff in index (if not a number, just hash)

type StripState struct {
	planJson         map[string]interface{}
	keyWhitelist     []string
	providers        []string
	deletes          []string
	hashedKeyParents []string // Parent keys within which to look for keys that need to be hashed

	moduleNames   []string
	variableNames []string
	resourceNames []string
	outputNames   []string
	indexNames    []string
}

// Helper function to produce hash digest of given string
func (S *StripState) Hash(value string) string {
	h := fnv.New64a()
	h.Write([]byte(value))

	// Prevent double hashing in recursion
	if strings.HasPrefix(value, "hash_") {
		return value
	}

	return fmt.Sprintf("hash_%v", h.Sum64())
}

func (S *StripState) HashOutputKeys(outputObject map[string]interface{}) map[string]interface{} {
	if outputValue, hasValue := outputObject["value"]; hasValue {
		if outputValue == nil {
			return outputObject
		}

		if reflect.TypeOf(outputValue).Kind() != reflect.Map {
			return outputObject // If value object is not a map -> no key hashing to be done
		}

		valueObject := outputObject["value"].(map[string]interface{})
		originalOutputValues := []string{}

		for valueKey, _ := range valueObject {
			originalOutputValues = append(originalOutputValues, valueKey)
		}

		for _, originalkey := range originalOutputValues {
			valueObject[S.Hash(originalkey)] = valueObject[originalkey]
			delete(valueObject, originalkey)
		}
	}

	return outputObject
}

// Function to handle name replacements in string values
func (S *StripState) CollectNames(inputMap map[string]interface{}) {
	for key, value := range inputMap {
		if value == nil {
			continue
		}

		// Handle provider config
		if key == "provider_config" {
			for _, providerValue := range value.(map[string]interface{}) {
				providerObject := providerValue.(map[string]interface{})
				S.providers = append(S.providers, providerObject["name"].(string))
			}
		}

		// Resource names
		if key == "resources" {
			for _, resourceValue := range value.([]interface{}) {
				resourceObject := resourceValue.(map[string]interface{})
				S.resourceNames = append(S.resourceNames, resourceObject["name"].(string))
			}
		}

		// Module names
		if key == "module_calls" {
			// Hash module keys and add new keys to module_calls object
			for moduleKey, _ := range value.(map[string]interface{}) {
				S.moduleNames = append(S.moduleNames, moduleKey)
			}

			// Delete old unhashed module keys
			for _, moduleName := range S.moduleNames {
				valueObject := value.(map[string]interface{})
				valueObject[S.Hash(moduleName)] = valueObject[moduleName]
				delete(value.(map[string]interface{}), moduleName)
			}
		}

		// Variable names
		if key == "variables" {
			for variableKey, _ := range value.(map[string]interface{}) {
				S.variableNames = append(S.variableNames, variableKey)
			}

			// Delete old unhashed module keys
			for _, variableName := range S.variableNames {
				valueObject := value.(map[string]interface{})

				// Check if value object actually has current name
				if _, hasOutput := valueObject[variableName]; hasOutput {
					valueObject[S.Hash(variableName)] = valueObject[variableName]
					delete(value.(map[string]interface{}), variableName)
				}
			}
		}

		// Output names
		if key == "outputs" {
			for outputKey, _ := range value.(map[string]interface{}) {
				S.outputNames = append(S.outputNames, outputKey)
			}

			// Delete old unhashed module keys
			for _, outputName := range S.outputNames {
				valueObject := value.(map[string]interface{})

				// Check if value object actually has current name
				if _, hasOutput := valueObject[outputName]; hasOutput {
					outputObject := valueObject[outputName].(map[string]interface{})
					outputObject = S.HashOutputKeys(outputObject)

					valueObject[S.Hash(outputName)] = outputObject
					delete(value.(map[string]interface{}), outputName)
				}
			}
		}

		// Proceed further down the recursion if value is type "map"
		if reflect.TypeOf(value).Kind() == reflect.Map {
			S.CollectNames(value.(map[string]interface{}))
		}
	}

}

// Function to handle name replacements in string values
func (S *StripState) ReplaceNames(inputValue string) string {
	replaceMatch := false
	inputParts := strings.Split(inputValue, ".")

	allNames := append(S.resourceNames, S.moduleNames...)
	allNames = append(allNames, S.variableNames...)

	for index, part := range inputParts {
		for _, name := range allNames {
			if part == name || strings.HasPrefix(part, name+"[") || strings.HasPrefix(part, name+":") {
				replaceMatch = true
				inputParts[index] = strings.ReplaceAll(part, name, S.Hash(name)) // Replace only name substring without altering or hashing index
			}
		}
	}

	if replaceMatch {
		return strings.Join(inputParts, ".")
	}

	return S.Hash(inputValue)
}

// Function to process all other value types
func (S *StripState) ProcessDefault(parentKey string, inputValue string) string {
	// Check if key is generally whitelisted
	for _, whitelistedKey := range S.keyWhitelist {
		if parentKey == whitelistedKey {
			return inputValue
		}
	}

	// Check if value is provider name
	for _, providerName := range S.providers {
		if inputValue == providerName {
			return inputValue
		}
	}

	whitelisted := false

	allNames := append(S.resourceNames, S.moduleNames...)
	allNames = append(allNames, S.variableNames...)

	for _, nameValue := range allNames {
		if strings.Contains(inputValue, nameValue) {
			whitelisted = true
			break
		}
	}

	if whitelisted {
		return S.ReplaceNames(inputValue)
	} else {
		return S.Hash(inputValue)
	}
}

// Function to recursively process slices
func (S *StripState) ProcessSlice(parentKey string, inputSlice []interface{}) {
	for index, value := range inputSlice {
		if value == nil {
			continue
		}

		valueType := reflect.TypeOf(value).Kind()

		switch valueType {
		case reflect.Map:
			S.ProcessMap(parentKey, value.(map[string]interface{}))
		case reflect.Slice:
			S.ProcessSlice(parentKey, value.([]interface{}))
		default:
			stringifiedValue := fmt.Sprintf("%v", value)
			processedValue := S.ProcessDefault(parentKey, stringifiedValue)

			// Converting potential numbers back to numbers if possible
			if auxiliary.IsNumeric(processedValue) {
				inputSlice[index], _ = strconv.Atoi(processedValue)
			} else {
				inputSlice[index] = processedValue
			}
		}
	}
}

// Function to recursively process maps
func (S *StripState) ProcessMap(parentKey string, inputMap map[string]interface{}) {
	// Delete unwanted keys
	for _, item := range S.deletes {
		delete(inputMap, item)
	}

	for key, value := range inputMap {
		if value == nil {
			continue
		}

		valueType := reflect.TypeOf(value).Kind()

		switch valueType {
		case reflect.Map:
			S.ProcessMap(key, value.(map[string]interface{}))
		case reflect.Slice:
			S.ProcessSlice(key, value.([]interface{}))
		default:
			stringifiedValue := fmt.Sprintf("%v", value)
			processedValue := S.ProcessDefault(key, stringifiedValue)

			if auxiliary.IsNumeric(processedValue) {
				inputMap[key], _ = strconv.Atoi(processedValue)
			} else {
				inputMap[key] = processedValue
			}
		}
	}
}

// Entrypoint
func (S *StripState) StripAndHash() error {
	functionName := "StripAndHashState"

	ux.PrintFormatted("⠿", []string{"blue"})
	ux.PrintFormatted(" Stripping Secrets", []string{"bold"})
	fmt.Println()
	fmt.Println()

	ux.PrintFormatted("→", []string{"blue"})
	fmt.Println(" We are stripping your plan state of secrets and hashing all values \n  to make it safe to share\n")

	stripSpinner := ux.NewSpinner("Stripping and hashing plan state", "Plan state stripped and hashed", "Stripping and hashing plan state failed", false)
	stripSpinner.Start()

	// Initialize relevant paths
	planPath := filepath.Join(auxiliary.PathInstance.WorkingPath, "pluralith.state.stripped")
	outPath := filepath.Join(auxiliary.PathInstance.WorkingPath, "pluralith.state.hashed")

	// Check if plan state exists -> if not, alert and return
	if _, existErr := os.Stat(planPath); existErr != nil {
		stripSpinner.Fail("No Pluralith plan state found")
		ux.PrintFormatted("→ Run pluralith plan again\n\n", []string{"red"})

		return nil
	}

	// Read plan state
	planBytes, readErr := os.ReadFile(planPath)
	if readErr != nil {
		stripSpinner.Fail("Failed to read plan state")
		return fmt.Errorf("could not read plan state -> %v: %w", functionName, readErr)
	}

	// Parse plan state
	parseErr := json.Unmarshal(planBytes, &S.planJson)
	if parseErr != nil {
		stripSpinner.Fail("Failed to parse plan state")
		return fmt.Errorf("could not parse plan state -> %v: %w", functionName, readErr)
	}

	S.CollectNames(S.planJson)
	S.providers = auxiliary.DeduplicateSlice(S.providers)
	S.resourceNames = auxiliary.DeduplicateSlice(S.resourceNames)
	S.moduleNames = auxiliary.DeduplicateSlice(S.moduleNames)
	S.variableNames = auxiliary.DeduplicateSlice(S.variableNames)

	S.keyWhitelist = []string{"type", "index", "provider_name", "terraform_version"}
	S.deletes = []string{"tags", "tags_all", "description", "source"}
	S.hashedKeyParents = []string{"value", "constant_value"}

	S.ProcessMap("", S.planJson)

	// Turn stripped and hashed state back into string
	planString, marshalErr := json.MarshalIndent(S.planJson, "", " ")
	if marshalErr != nil {
		stripSpinner.Fail("Failed to stringify stripped plan state")
		return fmt.Errorf("%v: %w", functionName, marshalErr)
	}

	// Write stripped and hashed state to file
	if writeErr := os.WriteFile(outPath, planString, 0700); writeErr != nil {
		stripSpinner.Fail("Failed to write stripped plan state")
		return fmt.Errorf("%v: %w", functionName, writeErr)
	}

	stripSpinner.Success()
	ux.PrintFormatted("→", []string{"blue"})
	fmt.Print(" Inspect it in the ")
	ux.PrintFormatted("pluralith.state.hashed", []string{"blue"})
	fmt.Println(" file")
	fmt.Println()

	return nil
}

var StripInstance = &StripState{}
