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
	"strings"
)

func hash(value string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(value))
	return h.Sum64()
}

// Helper function to check if value needs to be blacklisted
func CheckAndBlacklist(currentKey string, currentValue interface{}, keylist []string, blacklist *[]string) {
	// If any of the keys in the blacklist are present -> add value to blacklist
	for _, blackKey := range keylist {
		if currentKey == blackKey {
			// fmt.Println(outerKey, blackKey)
			stringified := fmt.Sprintf("%s", currentValue)
			*blacklist = append(*blacklist, stringified)
		}
	}
}

// Helper function to hash values (differentiates between array values and object key values)
func CheckAndHash(planJson map[string]interface{}, currentKey string, blacklist []string, index int) {
	var stringifiedValue string
	var blacklisted = false

	// Get value based on if array or not
	if index > -1 {
		slice := planJson[currentKey].([]interface{})
		stringifiedValue = fmt.Sprintf("%s", slice[index])
	} else {
		stringifiedValue = fmt.Sprintf("%s", planJson[currentKey])
	}

	// Check if blacklist contains value at current key
	for _, blackKey := range blacklist {
		if strings.Contains(blackKey, stringifiedValue) {
			blacklisted = true
			break
		}
	}

	// Set value based on if array or not
	if !blacklisted {
		if index > -1 {
			slice := planJson[currentKey].([]interface{})
			slice[index] = hash(stringifiedValue)
		} else {
			planJson[currentKey] = hash(stringifiedValue)
		}
	}
}

// Function to build a blacklist of values that should not be hashed
func BuildBlacklist(planJson map[string]interface{}, keylist []string, blacklist *[]string) {
	for key, value := range planJson {
		// Check if value at key is given
		if value != nil {
			outerValueType := reflect.TypeOf(value)

			// Switch between different data types
			switch outerValueType.Kind() {
			case reflect.Map:
				// If value is of type map -> Move on to next recursion level
				BuildBlacklist(value.(map[string]interface{}), keylist, blacklist)
			case reflect.Array, reflect.Slice:
				// If value is of type array or slice -> Loop through elements, if maps are found -> Move to next recursion level
				for _, sliceValue := range value.([]interface{}) {
					if reflect.TypeOf(sliceValue).Kind() == reflect.Map {
						BuildBlacklist(sliceValue.(map[string]interface{}), keylist, blacklist)
					} else {
						CheckAndBlacklist(key, sliceValue, keylist, blacklist)
					}
				}
			default:
				CheckAndBlacklist(key, value, keylist, blacklist)
			}
		}
	}
}

// Function to process plan state and strip all sensitive data, keeping structure intact for debugging
func ProcessState(planJson map[string]interface{}, blacklist []string) {
	for outerKey, outerValue := range planJson {
		// Check if value at key is given
		if outerValue != nil {
			outerValueType := reflect.TypeOf(outerValue)

			// Switch between different data types
			switch outerValueType.Kind() {
			case reflect.Map:
				// If value is of type map -> Move on to next recursion level
				ProcessState(outerValue.(map[string]interface{}), blacklist)
			case reflect.Array, reflect.Slice:
				// If value is of type array or slice -> Loop through elements, if maps are found -> Move to next recursion level
				for innerIndex, innerValue := range outerValue.([]interface{}) {
					if reflect.TypeOf(innerValue).Kind() == reflect.Map {
						ProcessState(innerValue.(map[string]interface{}), blacklist)
					} else {
						CheckAndHash(planJson, outerKey, blacklist, innerIndex)
					}
				}
			default:
				CheckAndHash(planJson, outerKey, blacklist, -1)
			}
		}
	}
}

func StripAndHash() error {
	functionName := "StripAndHashState"

	ux.PrintFormatted("⠿", []string{"blue"})
	ux.PrintFormatted(" Stripping Secrets", []string{"bold"})
	fmt.Println()
	fmt.Println()

	ux.PrintFormatted("→", []string{"blue"})
	fmt.Println(" We are stripping your plan state of secrets and hashing all values \n  to make it safe to share\n")

	stripSpinner := ux.NewSpinner("Stripping and hashing plan state", "Plan state stripped and hashed", "Stripping and hashing plan state failed", false)
	stripSpinner.Start()

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
	var planJson map[string]interface{}
	parseErr := json.Unmarshal(planBytes, &planJson)
	if parseErr != nil {
		stripSpinner.Fail("Failed to parse plan state")
		return fmt.Errorf("could not parse plan state -> %v: %w", functionName, readErr)
	}

	// Recursively collect exception values and build a blacklist
	keyBlacklist := []string{"address", "name", "type", "module_address", "index"}
	valueBlacklist := []string{}
	BuildBlacklist(planJson, keyBlacklist, &valueBlacklist)

	// Recursively process state
	ProcessState(planJson, valueBlacklist)

	// Turn stripped and hashed state back into string
	planString, marshalErr := json.MarshalIndent(planJson, "", " ")
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
