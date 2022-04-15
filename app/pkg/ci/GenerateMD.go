package ci

import (
	"fmt"
)

func GenerateMD(urls map[string]interface{}, changeActions map[string]interface{}) (string, error) {
	var comment string = ""

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

	return comment, nil
}
