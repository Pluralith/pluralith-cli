package ci

import "os"

func CheckDocker() bool {
	// Check if .dockerenv file exists
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
