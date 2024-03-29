package components

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/install"
	"runtime"
	"strings"
)

func GraphModule(silentCheck bool) {
	// Construct url
	url := "https://api.pluralith.com/v1/dist/download/cli/graphing"
	// url := "http://localhost:8080/v1/dist/download/cli/graphing"
	params := map[string]string{"os": runtime.GOOS, "arch": runtime.GOARCH}

	// Generate install path
	installPath := filepath.Join(auxiliary.StateInstance.BinPath, "pluralith-cli-graphing")
	if runtime.GOOS == "windows" {
		installPath = installPath + ".exe"
	}

	// Get current version
	var currentVersion string
	currentVersionByte, versionErr := exec.Command(installPath, "version").Output()
	if versionErr != nil {
		currentVersion = ""
	} else {
		currentVersion = strings.TrimSpace(string(currentVersionByte))
	}

	// Get Github release
	downloadUrl, shouldDownload, checkErr := install.GetGitHubRelease(url, params, currentVersion, silentCheck)
	if checkErr != nil {
		fmt.Println(checkErr)
	}

	// Handle download
	if shouldDownload {
		if downloadErr := install.DownloadGitHubRelease("Graph Module", downloadUrl, installPath); downloadErr != nil {
			fmt.Println(downloadErr)
		}
	}
}
