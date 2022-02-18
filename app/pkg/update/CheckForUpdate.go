package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
	"runtime"

	"github.com/hashicorp/go-version"
)

func CheckForUpdate() (string, bool, error) {
	functionName := "CheckForUpdate"

	checkSpinner := ux.NewSpinner("Checking for update", "You are on the latest version!\n", "Checking for update failed, try again!\n", false)

	request, _ := http.NewRequest("GET", "http://localhost:8080/v1/dist/download/cli", nil)
	request.Header.Add("Authorization", "Bearer "+os.Getenv("PluralithDevApiKey"))

	queryString := request.URL.Query()
	queryString.Add("os", runtime.GOOS)
	queryString.Add("arch", runtime.GOARCH)
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
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	parseErr := json.Unmarshal(bodyBytes, &bodyObject)
	if parseErr != nil {
		checkSpinner.Fail("Parsing request result failed")
		return "", false, fmt.Errorf("parsing request result failed -> %v: %w", functionName, responseErr)
	}

	versionData := bodyObject["data"].(map[string]interface{})

	currentVersion, _ := version.NewVersion(auxiliary.StateInstance.CLIVersion)
	latestVersion, _ := version.NewVersion(versionData["version"].(string))

	if currentVersion.LessThan(latestVersion) {
		// fmt.Printf("%s is less than %s\n", currentVersion, latestVersion)
		checkSpinner.Success("A new version is available!")

		ux.PrintFormatted("⠿ ", []string{"blue"})
		fmt.Print(auxiliary.StateInstance.CLIVersion + " → ")
		ux.PrintFormatted(latestVersion.Original()+"\n\n", []string{"blue", "bold"})

		return versionData["url"].(string), true, nil
	}

	checkSpinner.Success("You are on the latest version\n")
	return "", false, nil
}
