package ci

import (
	"os"
	"os/exec"
	"strings"
)

func GetBranch() string {
	// Get branch name from specific vendor env variables
	for _, vendor := range CIVendors {
		if vendor.Branch != "" {
			branch, found := os.LookupEnv(vendor.Branch)
			if found {
				return branch
			}
		}
	}

	// Attempt to get branch name via git if vendor specific approach fails
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchName, branchErr := branchCmd.Output()
	if branchErr == nil {
		return strings.TrimSpace(string(branchName))
	}

	return "none"
}
