package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pluralith/pkg/ux"

	"github.com/schollz/progressbar/v3"
)

func DownloadUpdate(url string) error {
	functionName := "DownloadUpdate"

	CLIPath, pathErr := os.Executable()
	if pathErr != nil {
		return fmt.Errorf("fetching latest version failed -> %v: %w", functionName, pathErr)
	}

	// Create bin file
	newFile, createErr := os.Create(CLIPath)
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
		progressbar.OptionSetDescription("Updating [light_blue]Pluralith[reset]"),
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

	ux.PrintFormatted("\n\n✔ Pluralith CLI updated!\n\n", []string{"green", "bold"})

	return nil
}
