package ci

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/cost"
	"pluralith/pkg/ux"
	"strconv"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func PostEvents(command string, tfArgs map[string]interface{}, costArgs map[string]interface{}, exportArgs map[string]interface{}) error {
	functionName := "PostEvents"

	caser := cases.Title(language.English)
	ux.PrintFormatted("\n→ ", []string{"blue", "bold"})
	ux.PrintFormatted(caser.String(command), []string{"white", "bold"})
	ux.PrintFormatted("\n  → ", []string{"blue"})
	fmt.Print("Terraform " + caser.String(command) + " Output: \n\n")

	var idStore = make(map[string]interface{})

	// Load cost cache
	costsByte, costsErr := os.ReadFile(filepath.Join(auxiliary.StateInstance.WorkingPath, ".pluralith", "pluralith.costs.json"))
	if costsErr != nil {
		return fmt.Errorf("loading infracost output failed -> %v: %w", functionName, costsErr)
	}
	costsMap := []cost.CostMap{}
	if parseErr := json.Unmarshal(costsByte, &costsMap); parseErr != nil {
		return fmt.Errorf("parsing infracost output failed -> %v: %w", functionName, parseErr)
	}

	// Generate resource cost dictionary
	resourceCosts := make(map[string]interface{})

	// Find costs for given resource
	for _, plan := range costsMap {
		for _, project := range plan.Projects {
			for _, resource := range project.Breakdown.Resources {
				costObject := ApplyEventCosts{}

				if resource.HourlyCost != nil {
					costObject.Hourly, _ = strconv.ParseFloat(resource.HourlyCost.(string), 64)
					costObject.Monthly, _ = strconv.ParseFloat(resource.MonthlyCost.(string), 64)
				}

				resourceCosts[resource.Name] = costObject
			}
		}
	}

	// Construct terraform args
	allArgs := []string{
		command,
		"-json",
		"-auto-approve",
	}

	// Construct arg slices for terraform
	for _, varValue := range tfArgs["var"].([]string) {
		allArgs = append(allArgs, "-var="+varValue)
	}

	for _, varFile := range tfArgs["var-file"].([]string) {
		allArgs = append(allArgs, "-var-file="+varFile)
	}

	// Run apply
	tfCmd := exec.Command("terraform", allArgs...)

	// Define sinks for std data
	var errorSink bytes.Buffer

	// Redirect command std data
	tfCmd.Stderr = &errorSink

	// Initiate standard output pipe
	outStream, outErr := tfCmd.StdoutPipe()
	if outErr != nil {
		return fmt.Errorf("streaming terraform output failed -> %v: %w", functionName, outErr)
	}

	// Run terraform command
	tfCmdErr := tfCmd.Start()
	if tfCmdErr != nil {
		return fmt.Errorf("running terraform apply failed -> %v: %w", functionName, tfCmdErr)
	}

	// Scan for command line updates
	applyScanner := bufio.NewScanner(outStream)
	applyScanner.Split(bufio.ScanLines)

	// While command line scan is running
	for applyScanner.Scan() {
		message := applyScanner.Text()

		// Parse terraform message
		parsedMessage := ApplyEvent{}
		parseErr := json.Unmarshal([]byte(message), &parsedMessage)
		if parseErr != nil {
			return fmt.Errorf("parsing terraform apply event failed -> %v: %w", functionName, parseErr)
		}

		// Print original Terraform apply event message
		if parsedMessage.Type != "version" && parsedMessage.Type != "planned_change" && parsedMessage.Type != "change_summary" {
			fmt.Println("    " + parsedMessage.Message)
		}

		// Filter for relevant apply events
		if parsedMessage.Type == "apply_start" || parsedMessage.Type == "apply_complete" || parsedMessage.Type == "apply_errored" {
			payload := make(map[string]interface{})

			// On delete: Save ID on apply start and append on apply complete (apply complete does not hold ID value anymore)
			if parsedMessage.Hook.Action == "delete" {
				if parsedMessage.Type == "apply_start" {
					idStore[parsedMessage.Hook.Resource.Addr] = parsedMessage.Hook.IDValue
				}
				if parsedMessage.Type == "apply_complete" {
					parsedMessage.Hook.IDValue = idStore[parsedMessage.Hook.Resource.Addr].(string)
				}
			}

			// Get resource cost from cost dictionary
			if resourceCosts[parsedMessage.Hook.Resource.Addr] != nil {
				parsedMessage.Hook.Resource.Costs = resourceCosts[parsedMessage.Hook.Resource.Addr].(ApplyEventCosts)
			}

			payload["projectId"] = auxiliary.StateInstance.PluralithConfig.ProjectId
			payload["runId"] = exportArgs["runId"]
			payload["event"] = parsedMessage

			// Encode payload
			payloadBytes, marshalErr := json.Marshal(payload)
			if marshalErr != nil {
				return fmt.Errorf("encoding terraform apply message failed -> %v: %w", functionName, marshalErr)
			}

			// Update resource
			// request, _ := http.NewRequest("POST", "https://api.pluralith.com/v1/resource/update", bytes.NewBuffer(messageBytes))
			request, _ := http.NewRequest("POST", "http://localhost:8080/v1/resource/update", bytes.NewBuffer(payloadBytes))
			request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)
			request.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			response, responseErr := client.Do(request)
			if responseErr != nil {
				return fmt.Errorf("sending resource update request failed -> %v: %w", functionName, responseErr)
			}

			// Parse response for file URLs
			responseBody, readErr := io.ReadAll(response.Body)
			if readErr != nil {
				return fmt.Errorf("reading resource update response failed -> %v: %w", functionName, readErr)
			}

			var bodyObject map[string]interface{}
			parseErr := json.Unmarshal(responseBody, &bodyObject)
			if parseErr != nil {
				return fmt.Errorf("parsing resource update response failed -> %v: %w", functionName, parseErr)
			}
		}
	}

	ux.PrintFormatted("\n✔ ", []string{"blue", "bold"})
	fmt.Print(caser.String(command) + " complete\n\n")

	return nil
}
