package auxiliary

import (
	"os"
)

func CheckEnvVars() bool {
	// Check for general matches
	for _, env := range GeneralEnvVars {
		if _, found := os.LookupEnv(env); found {
			return true
		}
	}

	// If no general match -> Check for vendor-specific matches
	for _, vendor := range CIVendors {
		for _, env := range vendor.Env {
			if _, found := os.LookupEnv(env); found {
				return true
			}
		}
	}

	// If no match found -> Not running in CI
	return false
}
