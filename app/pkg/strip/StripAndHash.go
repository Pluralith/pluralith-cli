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

type StripState struct {
	keyBlacklist   []string
	valueBlacklist []string
	planJson       map[string]interface{}
}

// Function to produce hash digest of given string
func (S *StripState) Hash(value string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(value))
	return h.Sum64()
}

// // Helper function to find and hash all resource names
// func ReplaceNames(planJson map[string]interface{}) {
// 	functionName := "ReplaceNames"

// 	resourceNames := []string{}
// 	for resource := range planJson["resources"] {

// 	}
// }

// Helper function to check if value needs to be blacklisted
func (S *StripState) CheckAndBlacklist(currentKey string, currentValue interface{}) {
	// If any of the keys in the blacklist are present -> add value to blacklist
	for _, blackKey := range S.keyBlacklist {
		if currentKey == blackKey {
			// fmt.Println(outerKey, blackKey)
			stringified := fmt.Sprintf("%s", currentValue) + "*"
			S.valueBlacklist = append(S.valueBlacklist, stringified)
		}
	}
}

// Helper function to hash values (differentiates between array values and object key values)
func (S *StripState) CheckAndHash(currentMap map[string]interface{}, currentKey string, index int) {
	var stringifiedValue string
	var blacklisted = false

	// Get value based on if array or not
	if index > -1 {
		slice := currentMap[currentKey].([]interface{})
		stringifiedValue = fmt.Sprintf("%s", slice[index])
	} else {
		stringifiedValue = fmt.Sprintf("%s", currentMap[currentKey])
	}

	// Check if blacklist contains value at current key
	for _, blackKey := range S.valueBlacklist {
		// Handle keys marked as prefixes (end with "*")
		if strings.HasSuffix(blackKey, "*") {
			noSuffixKey := strings.ReplaceAll(blackKey, "*", "")

			if strings.HasPrefix(stringifiedValue, noSuffixKey) {
				blacklisted = true
				break
			}
		}

		if strings.Contains(blackKey, stringifiedValue) {
			blacklisted = true
			break
		}
	}

	// Set value based on if array or not
	if !blacklisted {
		if index > -1 {
			slice := currentMap[currentKey].([]interface{})
			slice[index] = S.Hash(stringifiedValue)
		} else {
			currentMap[currentKey] = S.Hash(stringifiedValue)
		}
	}
}

// Function to build a blacklist of values that should not be hashed
func (S *StripState) BuildBlacklist(planMap map[string]interface{}) {
	for key, value := range planMap {
		// Check if value at key is given
		if value != nil {
			outerValueType := reflect.TypeOf(value)

			// Switch between different data types
			switch outerValueType.Kind() {
			case reflect.Map:
				currentMap := value.(map[string]interface{})
				// If value is of type map -> Move on to next recursion level
				S.BuildBlacklist(currentMap)
				// Checking for address
				// if _, hasAddress := currentMap["address"]; hasAddress {
				// 	if name, hasName := currentMap["name"]; hasName {
				// 		fmt.Println(name)
				// 	}
				// }

			case reflect.Array, reflect.Slice:
				// If value is of type array or slice -> Loop through elements, if maps are found -> Move to next recursion level
				for _, sliceValue := range value.([]interface{}) {
					if reflect.TypeOf(sliceValue).Kind() == reflect.Map {
						S.BuildBlacklist(sliceValue.(map[string]interface{}))
					} else {
						S.CheckAndBlacklist(key, sliceValue)
					}
				}
			default:
				S.CheckAndBlacklist(key, value)
			}
		}
	}
}

// Function to process plan state and strip all sensitive data, keeping structure intact for debugging
func (S *StripState) ProcessState(currentMap map[string]interface{}) {
	for outerKey, outerValue := range currentMap {
		// Check if value at key is given
		if outerValue != nil {
			outerValueType := reflect.TypeOf(outerValue)

			// Switch between different data types
			switch outerValueType.Kind() {
			case reflect.Map:
				// If value is of type map -> Move on to next recursion level
				S.ProcessState(outerValue.(map[string]interface{}))
			case reflect.Array, reflect.Slice:
				// If value is of type array or slice -> Loop through elements, if maps are found -> Move to next recursion level
				for innerIndex, innerValue := range outerValue.([]interface{}) {
					if reflect.TypeOf(innerValue).Kind() == reflect.Map {
						S.ProcessState(innerValue.(map[string]interface{}))
					} else {
						S.CheckAndHash(currentMap, outerKey, innerIndex)
					}
				}
			default:
				S.CheckAndHash(currentMap, outerKey, -1)
			}
		}
	}
}

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
	// var planJson map[string]interface{}
	parseErr := json.Unmarshal(planBytes, &S.planJson)
	if parseErr != nil {
		stripSpinner.Fail("Failed to parse plan state")
		return fmt.Errorf("could not parse plan state -> %v: %w", functionName, readErr)
	}

	// Replace resource names with hashes
	// var replaceErr error
	// planJson, replaceErr = ReplaceNames(planJson)
	// if replaceErr != nil {
	// 	stripSpinner.Fail("Failed to hash resource names")
	// 	return fmt.Errorf("could not hash resource names -> %v: %w", functionName, readErr)
	// }

	// Recursively collect exception values and build a blacklist
	S.keyBlacklist = []string{"address", "name", "type", "module_address", "index"}
	S.valueBlacklist = []string{"each.key", "count.index", "module.*", "var.*"}
	S.BuildBlacklist(S.planJson)

	// Deduplicate value blacklist
	S.valueBlacklist = auxiliary.DeduplicateSlice(S.valueBlacklist)

	// Recursively process state
	S.ProcessState(S.planJson)

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
