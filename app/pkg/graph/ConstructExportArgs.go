package graph

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"pluralith/pkg/auxiliary"
	"pluralith/pkg/ux"
	"strings"

	"github.com/spf13/pflag"
)

func ConstructExportArgs(flags *pflag.FlagSet, runAsCI bool) (map[string]interface{}, error) {
	functionName := "ConstructExportArgs"

	exportArgs := make(map[string]interface{})

	// Set variable values according to config export object if given
	exportArgs["Title"] = auxiliary.StateInstance.PluralithConfig.Export.Title
	exportArgs["Author"] = auxiliary.StateInstance.PluralithConfig.Export.Author
	exportArgs["Version"] = auxiliary.StateInstance.PluralithConfig.Export.Version

	// Get variable values that are relevant for input (override config values if flags are explicitly set)
	flagTitle, _ := flags.GetString("title")
	if flagTitle != "" {
		exportArgs["Title"] = flagTitle
	}
	flagAuthor, _ := flags.GetString("author")
	if flagAuthor != "" {
		exportArgs["Author"] = flagAuthor
	}
	flagVersion, _ := flags.GetString("version")
	if flagVersion != "" {
		exportArgs["Version"] = flagVersion
	}

	missingArguments := []string{}

	if exportArgs["Title"] == "" || exportArgs["Author"] == "" || exportArgs["Version"] == "" {
		ux.PrintFormatted("\n→", []string{"blue", "bold"})
		ux.PrintFormatted(" Details\n", []string{"white", "bold"})
	}

	// If not in CI -> Read all missing diagram values from stdin
	for key, _ := range exportArgs {
		if exportArgs[key] == "" {
			if runAsCI { // In CI -> Push key into missing arguments and fail after loop
				missingArguments = append(missingArguments, key)
			} else { // Locally -> Ask user for input
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
	}

	if len(missingArguments) != 0 {
		ux.PrintFormatted("✘", []string{"red", "bold"})
		fmt.Print(" Missing required arguments: ")
		for _, argument := range missingArguments {
			fmt.Print("\n  ⇢ ")
			ux.PrintFormatted(argument, []string{"red", "bold"})
		}
		fmt.Println()

		return exportArgs, errors.New("Handled")
	}

	// Get remaining diagram values that don't require potential user input
	exportArgs["OutDir"], _ = flags.GetString("out-dir")
	exportArgs["FileName"], _ = flags.GetString("file-name")
	exportArgs["SkipPlan"], _ = flags.GetBool("skip-plan")
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
