package install

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"

	"github.com/hashicorp/go-version"
)

func GetGitHubRelease(url string, params map[string]string, currentVersionString string) (string, bool, error) {
	functionName := "GetGitHubRelease"

	checkSpinner := ux.NewSpinner("Checking for update", "You are on the latest version!\n", "Checking for update failed, try again!\n", false)
	checkSpinner.Start()

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+auxiliary.StateInstance.APIKey)

	queryString := request.URL.Query()
	for paramKey, paramValue := range params {
		queryString.Add(paramKey, paramValue)
	}
	request.URL.RawQuery = queryString.Encode()

	// Execute get version request
	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil || response.StatusCode != 200 {
		checkSpinner.Fail("Fetching latest version failed")
		return "", false, fmt.Errorf("fetching latest version failed -> %v: %w", functionName, responseErr)
	}

	// Parse request body
	var bodyObject map[string]interface{}
	bodyBytes, _ := io.ReadAll(response.Body)
	parseErr := json.Unmarshal(bodyBytes, &bodyObject)
	if parseErr != nil {
		checkSpinner.Fail("Parsing request result failed")
		return "", false, fmt.Errorf("parsing response failed -> %v: %w", functionName, responseErr)
	}

	versionData := bodyObject["data"].(map[string]interface{})

	var currentVersion *version.Version
	var successMessage string

	// Handle non-existent version
	if len(currentVersionString) > 0 {
		currentVersion, _ = version.NewVersion(currentVersionString)
		successMessage = "A new version is available!"
	} else {
		currentVersion, _ = version.NewVersion("0.0.0")
		currentVersionString = "None"
		successMessage = "No graph module installed, found latest release"
	}

	latestVersion, _ := version.NewVersion(versionData["version"].(string))

	// Handle case if newer version is available
	if currentVersion.LessThan(latestVersion) {
		checkSpinner.Success(successMessage)

		ux.PrintFormatted("⠿ ", []string{"blue"})
		fmt.Print(currentVersionString + " → ")
		ux.PrintFormatted(latestVersion.Original()+"\n\n", []string{"blue", "bold"})

		return versionData["url"].(string), true, nil
	}

	// Else show that latest version is installed
	checkSpinner.Success("You are on the latest version")

	ux.PrintFormatted("⠿ ", []string{"bold", "blue"})
	fmt.Print("Version: ")
	ux.PrintFormatted(currentVersion.Original()+"\n\n", []string{"bold", "blue"})

	return "", false, nil
}
