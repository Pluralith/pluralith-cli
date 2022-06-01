package cost

import (
	"github.com/spf13/pflag"
)

func ConstructInfracostArgs(flags *pflag.FlagSet) map[string]interface{} {
	flagMap := make(map[string]interface{})

	flagMap["show-costs"], _ = flags.GetBool("show-costs")
	flagMap["cost-usage-file"], _ = flags.GetString("cost-usage-file")

	return flagMap
}
