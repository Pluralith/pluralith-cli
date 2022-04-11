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

func ConstructExportArgs(flags *pflag.FlagSet) (map[string]interface{}, error) {
	functionName := "ConstructExportArgs"

	exportArgs := make(map[string]interface{})

	// Get variable values that are relevant for input
	exportArgs["Title"], _ = flags.GetString("title")
	exportArgs["Author"], _ = flags.GetString("author")
	exportArgs["Version"], _ = flags.GetString("version")

	// Print UX head
	ux.PrintFormatted("⠿", []string{"blue", "bold"})
	if exportArgs["Title"] == "" && exportArgs["Author"] == "" && exportArgs["Version"] == "" {
		fmt.Println(" Exporting Diagram ⇢ Specify details below")
	} else {
		fmt.Println(" Exporting Diagram ⇢ Details taken from flags")
	}

	// Read all missing diagram values from stdin
	for key, _ := range exportArgs {
		if exportArgs[key] == "" {
			ux.PrintFormatted("  →", []string{"blue", "bold"})
			fmt.Printf(" %s: ", key)

			// Create scanner
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				exportArgs[key] = scanner.Text()
			}
			if scanErr := scanner.Err(); scanErr != nil {
				return exportArgs, fmt.Errorf("scanning input failed -> %v: %w", functionName, scanErr)
			}
		}
	}

	// Get remaining diagram values that don't require potential user input
	exportArgs["OutDir"], _ = flags.GetString("out-dir")
	exportArgs["FileName"], _ = flags.GetString("file-name")
	exportArgs["SkipPlan"], _ = flags.GetBool("skip-plan")
	exportArgs["GenerateMd"], _ = flags.GetBool("generate-md")
	exportArgs["ShowChanges"], _ = flags.GetBool("show-changes")
	exportArgs["ShowDrift"], _ = flags.GetBool("show-drift")

	// If no explicit output directory given -> Write to current working directory
	if exportArgs["OutDir"] == "" {
		exportArgs["OutDir"] = auxiliary.StateInstance.WorkingPath
	}

	// If no explicit file name given -> use title with spaces removed
	if exportArgs["FileName"] == "" { //, ".pdf", ""))
		exportArgs["FileName"] = strings.ReplaceAll(strings.ReplaceAll(exportArgs["Title"].(string), " ", ""), ".pdf", "") // Remove .pdf endings to avoid double file extensions
	}

	fmt.Println()

	return exportArgs, nil
}
