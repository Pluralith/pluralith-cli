package graph

import (
	"fmt"
	"path/filepath"
	"pluralith/pkg/auxiliary"
)

func GetAlgorithm() (map[string]interface{}, error) {
	functionName := "GetAlgorithm"

	// Define path to graphing cache
	cachePath := filepath.Join(auxiliary.StateInstance.PluralithPath, "pluralithCache.json")
	// Get installed graphing version
	installedVersion, versionErr := FetchInstalledAlgorithm(cachePath)
	if versionErr != nil {
		// If no version installed -> fetch latest
		latestAlgorithm, fetchErr := FetchLatestAlgorithm(cachePath)
		if fetchErr != nil {
			return latestAlgorithm, fmt.Errorf("failed to fetch latest graphing data -> %v: %w", functionName, fetchErr)
		}

		return latestAlgorithm, nil
	}

	// Get latest production version string
	latestVersion, versionErr := FetchLatestVersion()
	if versionErr != nil {
		fmt.Println(fmt.Errorf("failed to fetch latest graphing version -> %v: %w", functionName, versionErr))
	}

	// Compare installed and latest version
	if installedVersion["version"] != latestVersion {
		// If latest version not equal to installed version -> fetch latest
		latestAlgorithm, fetchErr := FetchLatestAlgorithm(cachePath)
		if fetchErr != nil {
			return latestAlgorithm, fmt.Errorf("%v: %w", functionName, fetchErr)
		}

		return latestAlgorithm, nil
	}

	fmt.Println("Algorithm version " + installedVersion["version"].(string) + " up to date")

	return installedVersion, nil
}
