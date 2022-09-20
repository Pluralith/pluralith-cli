package graph

import (
	"pluralith/pkg/auxiliary"

	"github.com/spf13/pflag"
)

func ConstructExportArgs(flags *pflag.FlagSet) map[string]interface{} {
	flagMap := make(map[string]interface{})

	flagMap["title"], _ = flags.GetString("title")
	flagMap["author"], _ = flags.GetString("author")
	flagMap["version"], _ = flags.GetString("version")
	flagMap["out-dir"], _ = flags.GetString("out-dir")
	flagMap["file-name"], _ = flags.GetString("file-name")
	flagMap["show-changes"], _ = flags.GetBool("show-changes")
	flagMap["show-drift"], _ = flags.GetBool("show-drift")
	flagMap["show-costs"], _ = flags.GetBool("show-costs")
	flagMap["export-pdf"], _ = flags.GetBool("export-pdf")
	flagMap["post-apply"], _ = flags.GetBool("post-apply")
	flagMap["branch"] = auxiliary.StateInstance.Branch

	return flagMap
}
