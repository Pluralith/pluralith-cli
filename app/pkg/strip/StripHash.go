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
	"regexp"
	"strconv"
	"strings"
)

// Entrypoint
func StripAndHash() error {
	functionName := "StripAndHashState"

	ux.PrintFormatted("⠿", []string{"blue"})
	ux.PrintFormatted(" Stripping Secrets", []string{"bold"})
	fmt.Println()
	fmt.Println()

	ux.PrintFormatted("→", []string{"blue"})
	fmt.Println(" We are stripping your plan state of secrets and hashing all values \n  to make it safe to share")

	stripSpinner := ux.NewSpinner("Stripping and hashing plan state", "Plan state stripped and hashed", "Stripping and hashing plan state failed", false)
	stripSpinner.Start()

	var planJson interface{} // State

	// Initialize relevant paths
	planPath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.state.json")
	outPath := filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.state.hashed")

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
	parseErr := json.Unmarshal(planBytes, &planJson)
	if parseErr != nil {
		stripSpinner.Fail("Failed to parse plan state")
		return fmt.Errorf("could not parse plan state -> %v: %w", functionName, readErr)
	}

	// Strip state
	planJson = StripJson(planJson)

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

// Strip state recursively
func StripJson(obj interface{}) interface{} {

	if obj == nil {
		return nil
	}

	objType := reflect.TypeOf(obj).Kind()
	if objType == reflect.String { // obj is string
		return HashString(obj.(string))
	} else if objType == reflect.Slice { // obj is array
		for i := 0; i < len(obj.([]interface{})); i++ {
			stripResult := StripJson((obj.([]interface{}))[i])
			if stripResult != nil {
				obj.([]interface{})[i] = stripResult
			}
		}
		return obj
	} else if objType == reflect.Map { // obj is map
		newMap := make(map[string]interface{})
		for k, value := range obj.(map[string]interface{}) {
			if value == nil {
				// TODO
			} else if reflect.TypeOf(value).Kind() == reflect.Bool {
				obj.(map[string]interface{})[k] = Hash(strconv.FormatBool(obj.(map[string]interface{})[k].(bool)))
			} else if reflect.TypeOf(value).Kind() == reflect.Int {
				// TODO
			} else if reflect.TypeOf(value).Kind() == reflect.String && value.(string) == "" {
				obj.(map[string]interface{})[k] = ""
			} else {
				stripResult := StripJson(value)
				if stripResult != nil {
					obj.(map[string]interface{})[k] = stripResult
				}
			}

			stripResult := HashString(k)
			if stripResult != k {
				newMap[stripResult] = obj.(map[string]interface{})[k]
			} else {
				newMap[k] = obj.(map[string]interface{})[k]
			}
		}

		obj = newMap
		return obj
	}

	return obj
}

// Hash a string using som e rules
func HashString(value string) string {

	if strings.Contains(value, "\n") || strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") || strings.HasSuffix(value, "}") {
		return Hash(value)
	}
	splitString := strings.Split(value, ".")

	for index1 := 0; index1 < len(splitString); index1++ {
		splitPart := []string{}
		if !strings.Contains(splitString[index1], "[") {
			splitPart = strings.Split(splitString[index1], "/")
		} else {
			splitPart = append(splitPart, splitString[index1])
		}

		for index := 0; index < len(splitPart); index++ {
			if splitPart[index] == "" {
				continue
			} else if strings.Contains(splitPart[index], "[") {
				r, _ := regexp.Compile("\\[((([^\\]])|(\"))+)\\]")
				regexMatch := r.FindStringSubmatch(splitPart[index])
				if regexMatch != nil {
					bracketContent := strings.ReplaceAll(regexMatch[1], "[", "")
					bracketContent = strings.ReplaceAll(regexMatch[1], "]", "")
					if _, err := strconv.Atoi(bracketContent); err != nil {
						quotations := false
						bracketContentBackup := bracketContent
						if strings.HasSuffix(bracketContent, "\"") && strings.HasPrefix(bracketContent, "\"") {
							bracketContent = strings.TrimPrefix(bracketContent, "\"")
							bracketContent = strings.TrimSuffix(bracketContent, "\"")
							quotations = true
						}
						bracketContent = Hash(bracketContent)
						if quotations {
							bracketContent = "\"" + bracketContent + "\""
						}

						splitPart[index] = strings.ReplaceAll(splitPart[index], "["+bracketContentBackup+"]", "["+bracketContent+"]")
					}

					splitPartPart := strings.Split(splitPart[index], "[")
					for index2 := 0; index2 < len(splitPartPart); index2++ {
						if splitPartPart[index2] == "" {
							continue
						} else if !strings.Contains(splitPartPart[index2], "]") {
							splitPartPart[index2] = Hash(splitPartPart[index2])
						} else {
							splitPartPart[index2] = "[" + splitPartPart[index2]
						}
					}
					hashedPartPart := ""

					for index2 := 0; index2 < len(splitPartPart); index2++ {
						hashedPartPart += splitPartPart[index2]
					}
					splitPart[index] = hashedPartPart
				}
			} else if !contains(GetStripBlacklist(), splitPart[index]) {
				if _, err := strconv.Atoi(splitPart[index]); err != nil {
					splitPart[index] = Hash(splitPart[index])
				}
			}
		}

		hashedPart := ""

		for index := 0; index < len(splitPart); index++ {
			hashedPart += splitPart[index]
			if index < len(splitPart)-1 {
				hashedPart += "/"
			}
		}
		splitString[index1] = hashedPart
	}

	hashedString := ""

	for index := 0; index < len(splitString); index++ {
		hashedString += splitString[index]
		if index < len(splitString)-1 {
			hashedString += "."
		}
	}

	return hashedString
}

// Helper function to produce hash digest of given string
func Hash(value string) string {
	h := fnv.New64a()
	h.Write([]byte(value))
	return fmt.Sprintf("hash_%v", h.Sum64())
}

// Check if string contains substring
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
