package graph

import (
	"bufio"
	"fmt"
	"os"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
	"strings"

	"github.com/spf13/pflag"
)

func GetDiagramValues(flags *pflag.FlagSet) (map[string]interface{}, error) {
	functionName := "GetDiagramValues"

	diagramValues := make(map[string]interface{})

	// Get variable values that are relevant for input
	diagramValues["Title"], _ = flags.GetString("title")
	diagramValues["Author"], _ = flags.GetString("author")
	diagramValues["Version"], _ = flags.GetString("version")

	// Print UX head
	ux.PrintFormatted("⠿", []string{"blue", "bold"})
	if diagramValues["Title"] == "" && diagramValues["Author"] == "" && diagramValues["Version"] == "" {
		fmt.Println(" Exporting Diagram ⇢ Specify details below")
	} else {
		fmt.Println(" Exporting Diagram ⇢ Details taken from flags")
	}

	// Read all missing diagram values from stdin
	for key, _ := range diagramValues {
		if diagramValues[key] == "" {
			ux.PrintFormatted("  →", []string{"blue", "bold"})
			fmt.Printf(" %s: ", key)

			// Create scanner
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				diagramValues[key] = scanner.Text()
			}
			if scanErr := scanner.Err(); scanErr != nil {
				return diagramValues, fmt.Errorf("scanning input failed -> %v: %w", functionName, scanErr)
			}
		}
	}

	// Get remaining diagram values that don't require potential user input
	diagramValues["OutDir"], _ = flags.GetString("out-dir")
	diagramValues["FileName"], _ = flags.GetString("file-name")
	diagramValues["SkipPlan"], _ = flags.GetBool("skip-plan")
	diagramValues["GenerateMd"], _ = flags.GetBool("generate-md")
	diagramValues["ShowChanges"], _ = flags.GetBool("show-changes")

	// If no explicit output directory given -> Write to current working directory
	if diagramValues["OutDir"] == "" {
		diagramValues["OutDir"] = auxiliary.StateInstance.WorkingPath
	}

	// If no explicit file name given -> use title with spaces removed
	if diagramValues["FileName"] == "" { //, ".pdf", ""))
		diagramValues["FileName"] = strings.ReplaceAll(strings.ReplaceAll(diagramValues["Title"].(string), " ", ""), ".pdf", "") // Remove .pdf endings to avoid double file extensions
	}

	fmt.Println()

	return diagramValues, nil
}
