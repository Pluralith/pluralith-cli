package terraform

import (
	"pluralith/pkg/auxiliary"

	"github.com/spf13/pflag"
)

func ConstructTerraformArgs(flags *pflag.FlagSet) map[string]interface{} {
	flagMap := make(map[string]interface{})

	flagMap["var"], _ = flags.GetStringArray("var")
	flagMap["var-file"], _ = flags.GetStringArray("var-file")
	flagMap["plan-file"], _ = flags.GetString("plan-file")
	flagMap["plan-file-json"], _ = flags.GetString("plan-file-json")

	// If no plan file of any form is provided -> Make sure to check for terraform init
	if flagMap["plan-file"] == "" && flagMap["plan-file-json"] == "" {
		auxiliary.StateInstance.CheckTerraformInit()
	}

	return flagMap
}
