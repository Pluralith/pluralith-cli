package install

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pluralith/pkg/ux"

	"github.com/schollz/progressbar/v3"
)

func DownloadGitHubRelease(name string, url string, installPath string) error {
	functionName := "DownloadGitHubRelease"

	// Create bin file
	newFile, createErr := os.Create(installPath)
	if createErr != nil {
		newFile.Close()
		return fmt.Errorf("fetching latest version failed -> %v: %w", functionName, createErr)
	}

	defer newFile.Close()

	// Get latest version
	response, getErr := http.Get(url)
	if getErr != nil {
		response.Body.Close()
		return fmt.Errorf("getting latest version failed -> %v: %w", functionName, getErr)
	}

	defer response.Body.Close()

	// Instantiate download bar
	downloadBar := progressbar.NewOptions64(response.ContentLength,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("Installing Latest [light_blue]"+name+"[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer: "[light_blue]█[reset]",
			// SaucerHead:    "[green]>[reset]",
			SaucerPadding: "[dark_gray]█[reset]",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	// Write to file and progress bar
	if _, writeErr := io.Copy(io.MultiWriter(newFile, downloadBar), response.Body); writeErr != nil {
		return fmt.Errorf("downloading latest version failed -> %v: %w", functionName, getErr)
	}

	chmodErr := os.Chmod(installPath, 0700)
	if chmodErr != nil {
		return fmt.Errorf("chmod on executable failed -> %v: %w", functionName, getErr)
	}

	ux.PrintFormatted("\n\n✔ "+name+" updated!\n\n", []string{"green", "bold"})

	return nil
}
