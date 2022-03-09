package ci

import (
	"fmt"
)

func GeneratePRComment(urls map[string]string) (string, error) {
	var comment string = ""

	// Generate Head
	comment += "## 📝 Pluralith Diagram Generated  \n\n"

	// Generate Body
	comment += fmt.Sprintf("![Pluralith Diagram](%s)  \n", urls["PNG"])
	comment += fmt.Sprintf("### View the [PDF](%s)  \n", urls["PDF"])

	return comment, nil
}
