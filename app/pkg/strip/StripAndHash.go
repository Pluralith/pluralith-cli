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
	"sort"
	"strings"
)

// Need to hash:
// - Resource names inside strings
// - - - > If map has address and name -> get name
// - Variable names inside strings
// - - - > If current key is "variables" -> get all map keys
// - Module names inside strings
// - - - > If current key is "module_calls" -> get all map keys

// CHECK OUTCOME BEFORE PROCEEDING

// - Every other string / number value
// - Resource names in keys
// - Variable names in keys
// - Handle arrays
// - Handle objects

// Exceptions:
// - Provider names
type StripState struct {
	planJson  map[string]interface{}
	names     []string
	whitelist []string
	deletes   []string
}

// Helper function to produce hash digest of given string
func (S *StripState) Hash(value string) string {
	h := fnv.New64a()
	h.Write([]byte(value))
	return fmt.Sprintf("hash_%v", h.Sum64())
}

// Helper function to handle map in recursion
func (S *StripState) HandleMap(inputKey string, inputMap map[string]interface{}) {
	// Handle general value case
	if _, hasAddress := inputMap["address"]; hasAddress {
		if name, hasName := inputMap["name"]; hasName {
			S.names = append(S.names, name.(string))
		}
	}

	// Handle special key case for module names
	if inputKey == "module_calls" {
		for moduleKey, _ := range inputMap {
			S.names = append(S.names, moduleKey)
		}
	}

	// Handle special key case for variable names
	if inputKey == "variables" {
		for variableKey, _ := range inputMap {
			S.names = append(S.names, variableKey)
		}
	}

	// Remove provider names from hash name list
	if inputKey == "provider_config" {
		for _, providerObject := range inputMap {
			mapConversion := providerObject.(map[string]interface{})
			S.whitelist = append(S.whitelist, mapConversion["name"].(string))
		}
	}
}

// Recursive function to get all resource, variable and module names
func (S *StripState) FetchNames(inputMap map[string]interface{}) {
	for key, value := range inputMap {
		if value != nil {
			valueType := reflect.TypeOf(value).Kind()

			switch valueType {
			case reflect.Map:
				S.HandleMap(key, value.(map[string]interface{})) // Get names variables or modules
				S.FetchNames(value.(map[string]interface{}))
			case reflect.Slice:
				for _, item := range value.([]interface{}) {
					if reflect.TypeOf(item).Kind() == reflect.Map {
						S.HandleMap("", item.(map[string]interface{}))
					}
				}
			}
		}
	}
}

// Function to check value for conditions and hash accordingly
func (S *StripState) CheckAndHash(inputMap map[string]interface{}, key string, index int) {
	whitelisted := false
	stringifiedValue := fmt.Sprintf("%v", inputMap[key])

	if index > -1 {
		stringifiedValue = fmt.Sprintf("%v", inputMap[key].([]interface{})[index])
	}

	// Handle whitelist items
	for _, item := range S.whitelist {
		// Handle values
		if stringifiedValue == item {
			whitelisted = true
			break
		}

		// Handle keys
		if key == item {
			whitelisted = true
			break
		}
	}

	// Handle names
	for _, name := range S.names {
		if strings.Contains(stringifiedValue, name) {
			whitelisted = true
			// fmt.Println(stringifiedValue, "-> contains ->", name)
			if index > -1 {
				keyValue := inputMap[key].([]interface{})
				keyValue[index] = strings.ReplaceAll(stringifiedValue, name, S.Hash(name))
				stringifiedValue = keyValue[index].(string)
			} else {
				inputMap[key] = strings.ReplaceAll(stringifiedValue, name, S.Hash(name))
				stringifiedValue = inputMap[key].(string)
			}

			// Add hashed names to whitelist to prevent them from being hashed again
			S.names = append(S.names, S.Hash(name))
		}
	}

	// Handle remaining values
	if !whitelisted {
		if index > -1 {
			keyValue := inputMap[key].([]interface{})
			keyValue[index] = S.Hash(stringifiedValue)
		} else {
			inputMap[key] = S.Hash(stringifiedValue)
		}
	}
}

// Recursive function to hash name appearances and other values in plan json
func (S *StripState) HashAppearances(inputMap map[string]interface{}) {
	for key, value := range inputMap {
		if value != nil {
			valueType := reflect.TypeOf(value).Kind()

			// Delete unwanted keys
			for _, item := range S.deletes {
				delete(inputMap, item)
			}

			switch valueType {
			case reflect.Map:
				S.HashAppearances(value.(map[string]interface{}))
			case reflect.Slice:
				for index, item := range value.([]interface{}) {
					if reflect.TypeOf(item).Kind() == reflect.Map {
						S.HashAppearances(item.(map[string]interface{}))
					} else {
						S.CheckAndHash(inputMap, key, index)
					}
				}
			default:
				S.CheckAndHash(inputMap, key, -1)
			}
		}

		// Handle key replacement
		for _, name := range S.names {
			if key == name {
				inputMap[S.Hash(key)] = value
				delete(inputMap, key)
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

	S.whitelist = []string{"address", "type", "module_address", "index", "provider_name", "each.key", "count.index"}
	S.deletes = []string{"tags", "tags_all", "description", "source"}

	// Fetch names
	S.FetchNames(S.planJson)
	S.names = auxiliary.DeduplicateSlice(S.names)

	// Sort name list by length to avoid erroneous substring matches
	sort.Slice(S.names, func(i, j int) bool {
		return len(S.names[i]) > len(S.names[j])
	})

	// Hash values
	S.HashAppearances(S.planJson)

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
