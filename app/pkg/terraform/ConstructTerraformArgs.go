package terraform

import (
	"github.com/spf13/pflag"
)

func ConstructTerraformArgs(flags *pflag.FlagSet) map[string]interface{} {
	flagMap := make(map[string]interface{})

	flagMap["var"], _ = flags.GetStringArray("var")
	flagMap["var-file"], _ = flags.GetStringArray("var-file")
	flagMap["plan-file"], _ = flags.GetString("plan-file")
	flagMap["plan-file-json"], _ = flags.GetString("plan-file-json")

	return flagMap
}
