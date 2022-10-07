package ci

import (
	"fmt"
	"os"
	"path/filepath"
	"pluralith/pkg/auxiliary"
)

func GenerateComment(runCache map[string]interface{}) error {
	functionName := "GenerateComment"

	comment := ""
	urls := runCache["urls"].(map[string]interface{})
	changeActions := runCache["changes"].(map[string]interface{})

	// Generate Head
	comment += "[![Pluralith GitHub Badge](https://user-images.githubusercontent.com/25454503/158065018-55796de7-60a8-4c91-8aa4-3f53cd3c253f.svg)](https://www.pluralith.com)\n\n"
	comment += "## Diagram Generated\n\n"

	// Generate Body
	// Diagram Section
	comment += "→ **`Click the image to view this run in the Pluralith Dashboard`**\n\n"
	comment += fmt.Sprintf("[![Pluralith Diagram](%s)](%s)\n\n", urls["thumbnailURL"], urls["pluralithURL"])

	// Changes Section
	comment += "## Changes\n\n"
	comment += fmt.Sprintf("| **Created** | **Updated** | **Destroyed** | **Recreated** | **Drifted** | **Unchanged** |\n|-------------|-------------|-------------|---------------|---------------|---------------|\n| 🟢 **`+ %v`** | 🟠 **`~ %v`** | 🔴 **`- %v`**   | 🔵 **`@ %v`**   | 🟣 **`# %v`**   | ⚪ **`# %v`**   |", changeActions["create"], changeActions["update"], changeActions["delete"], changeActions["deletecreate"], changeActions["drift"], changeActions["no-op"])

	// Write markdown to file system for usage by pipeline
	commentPath := filepath.Join(auxiliary.StateInstance.WorkingPath, "comment.md")
	if writeErr := os.WriteFile(commentPath, []byte(comment), 0700); writeErr != nil {
		return fmt.Errorf("writing PR comment markdown to filesystem failed -> %v: %w", functionName, writeErr)
	}

	return nil
}
