package terraform

import (
	"github.com/spf13/pflag"
)

func ConstructTerraformArgs(flags *pflag.FlagSet) map[string]interface{} {
	flagMap := make(map[string]interface{})

	flagMap["var"], _ = flags.GetStringSlice("var")
	flagMap["var-file"], _ = flags.GetStringSlice("var-file")

	return flagMap
}
