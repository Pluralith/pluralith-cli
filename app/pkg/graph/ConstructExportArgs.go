package graph

import (
	"fmt"
	"pluralith/pkg/auxiliary"

	"github.com/spf13/pflag"
)

func ConstructExportArgs(flags *pflag.FlagSet) map[string]interface{} {
	functionName := "ConstructExportArgs"

	flagMap := make(map[string]interface{})

	flagMap["local-only"], _ = flags.GetBool("local-only")
	flagMap["title"], _ = flags.GetString("title")
	flagMap["author"], _ = flags.GetString("author")
	flagMap["version"], _ = flags.GetString("version")
	flagMap["out-dir"], _ = flags.GetString("out-dir")
	flagMap["file-name"], _ = flags.GetString("file-name")
	flagMap["show-changes"], _ = flags.GetBool("show-changes")
	flagMap["show-drift"], _ = flags.GetBool("show-drift")
	flagMap["show-costs"], _ = flags.GetBool("show-costs")
	flagMap["export-pdf"], _ = flags.GetBool("export-pdf")
	flagMap["sync-to-backend"], _ = flags.GetBool("sync-to-backend")
	flagMap["post-apply"], _ = flags.GetBool("post-apply")
	flagMap["branch"] = auxiliary.StateInstance.Branch

	// Load custom config file if specified by user
	customConfigPath, _ := flags.GetString("config-file")
	if getConfigErr := auxiliary.StateInstance.GetConfig(customConfigPath); getConfigErr != nil {
		fmt.Println(fmt.Errorf("fetching pluralith config failed -> %v: %w", functionName, getConfigErr))
	}

	// Get title from config if not passed as flag
	if flagMap["title"] == "" {
		flagMap["title"] = auxiliary.StateInstance.PluralithConfig.Config.Title
	}

	// Get version from config if not passed as flag
	if flagMap["version"] == "" {
		flagMap["version"] = auxiliary.StateInstance.PluralithConfig.Config.Version
	}

	return flagMap
}
