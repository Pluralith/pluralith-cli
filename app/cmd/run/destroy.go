package run

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ci"
	"pluralith/pkg/install/components"
	"pluralith/pkg/ux"

	"github.com/spf13/cobra"
)

var RunDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy and push updates to Pluralith",
	Long:  `Run terraform destroy and push updates to Pluralith`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print UX head
		ux.PrintFormatted("⠿", []string{"blue", "bold"})
		fmt.Println(" Initiating Apply Run ⇢ Posting To Pluralith Dashboard")

		var idStore = make(map[string]interface{})

		tfArgs, costArgs, exportArgs, preErr := ci.PreRun(cmd.Flags())
		if preErr != nil {
			fmt.Println(preErr, costArgs, exportArgs)
		}

		// Check if graph module installed, if not -> install
		_, versionErr := exec.Command(filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing"), "version").Output()
		if versionErr != nil {
			components.GraphModule()
		}

		// Construct terraform args
		allArgs := []string{
			"destroy",
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
			fmt.Println(outErr)
			// return fmt.Errorf("%v: %w", functionName, outErr)
		}

		// Run terraform command
		tfCmdErr := tfCmd.Start()
		if tfCmdErr != nil {
			fmt.Println(tfCmdErr)
		}

		// Scan for command line updates
		applyScanner := bufio.NewScanner(outStream)
		applyScanner.Split(bufio.ScanLines)

		fmt.Println(auxiliary.StateInstance.PluralithConfig.ProjectId, exportArgs["runId"])

		// eventLog := [][]string{}
		// eventPadding := 0

		// errorCount := 0
		// errorPrint := color.New(color.Bold, color.FgHiRed)

		// successCount := 0
		// successMode := "Created"
		// successPrint := color.New(color.Bold, color.FgHiGreen)

		// if command == "destroy" {
		// 	successMode = "Destroyed"
		// 	successPrint = color.New(color.Bold, color.FgHiBlue)
		// }

		// Deactivate cursor
		// fmt.Print("\033[?25l")

		// ux.PrintFormatted("  → ", []string{"bold", "blue"})
		// fmt.Printf("Running → %s %s / %s Errored", successPrint.Sprint(strconv.Itoa(successCount)), successMode, errorPrint.Sprint(strconv.Itoa(errorCount)))

		// While command line scan is running
		for applyScanner.Scan() {
			message := applyScanner.Text()

			// Parse terraform message
			parsedMessage := ci.ApplyEvent{}
			parseErr := json.Unmarshal([]byte(message), &parsedMessage)
			if parseErr != nil {
				fmt.Println(parseErr)
				return
			}

			// Filter for relevant apply events
			if parsedMessage.Type == "apply_start" || parsedMessage.Type == "apply_complete" || parsedMessage.Type == "apply_errored" {
				payload := make(map[string]interface{})
				address := parsedMessage.Hook.Resource.Addr
				action := parsedMessage.Hook.Action

				// On delete: Save ID on apply start and append on apply complete (apply complete does not hold ID value anymore)
				if action == "delete" {
					if parsedMessage.Type == "apply_start" {
						idStore[address] = parsedMessage.Hook.IDValue
					}
					if parsedMessage.Type == "apply_complete" {
						parsedMessage.Hook.IDValue = idStore[address].(string)
					}
				}

				payload["projectId"] = auxiliary.StateInstance.PluralithConfig.ProjectId
				payload["runId"] = exportArgs["runId"]
				payload["event"] = parsedMessage

				payloadBytes, marshalErr := json.Marshal(payload)
				if marshalErr != nil {
					fmt.Println(marshalErr)
					return
				}

				// request, _ := http.NewRequest("POST", "https://api.pluralith.com/v1/resource/update", bytes.NewBuffer(messageBytes))
				request, _ := http.NewRequest("POST", "http://localhost:8080/v1/resource/update", bytes.NewBuffer(payloadBytes))
				request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)
				request.Header.Add("Content-Type", "application/json")

				client := &http.Client{}
				response, responseErr := client.Do(request)

				if responseErr != nil {
					fmt.Println(responseErr)
					return
				}

				// Parse response for file URLs
				responseBody, readErr := ioutil.ReadAll(response.Body)
				if readErr != nil {
					fmt.Println(readErr)
					return
				}

				var bodyObject map[string]interface{}
				parseErr := json.Unmarshal(responseBody, &bodyObject)
				if parseErr != nil {
					fmt.Println(parseErr)
					return
				}

				// fmt.Println(responseBody)

			}

			// fmt.Println(parsedMessage["@message"])
			fmt.Println(message)
		}
	},
}

func init() {}
