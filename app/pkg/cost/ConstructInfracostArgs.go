package cost

import (
	"fmt"

	"github.com/spf13/pflag"
)

func ConstructInfracostArgs(flags *pflag.FlagSet) (map[string]interface{}, error) {
	flagMap := make(map[string]interface{})

	flagMap["show-costs"], _ = flags.GetBool("show-costs")
	flagMap["cost-usage-file"], _ = flags.GetString("cost-usage-file")
	flagMap["cost-mode"], _ = flags.GetString("cost-mode")
	flagMap["cost-period"], _ = flags.GetString("cost-period")

	if flagMap["cost-mode"] != "delta" && flagMap["cost-mode"] != "total" {
		return flagMap, fmt.Errorf("Invalid value for --cost-mode. Can only be 'delta' or 'total'")
	}

	fmt.Println(flagMap["cost-period"])

	if flagMap["cost-period"] != "hour" && flagMap["cost-period"] != "month" {
		return flagMap, fmt.Errorf("Invalid value for --cost-period. Can only be 'hour' or 'month'")
	}

	return flagMap, nil
}
