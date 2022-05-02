package auxiliary

import (
	"strings"
)

// - - - Code to manually parse flags - - -
// (Need to do it manually because cobra won't allow for unknown flags to pass through to terraform commands)

// Function to parse flags
func ParseArgs(args []string, pluralithArgs map[string]string) []string {
	// Instantiating clean slice to collect args
	var parsedArgs []string
	parsedArgMap := make(map[string]string)

	// Cleaning up potential '=' in args
	for _, arg := range args {
		splitArg := strings.Split(arg, "=")
		parsedArgs = append(parsedArgs, splitArg...)
	}

	totalArgs := len(parsedArgs) - 1 // Getting cleanArgs length to use in loop

	// Creating arg map
	for index, arg := range parsedArgs {
		// Determining if current value is an arg or a value by checking for "-"
		if strings.HasPrefix(arg, "-") {
			// Initializing a few variables
			var value string

			currentArg := arg[1:] // Clearing "-" prefix from args

			// If there is a value in arg slice beyond this one and it isn't another flag (doesn't start with '-') -> treat it as value for current flag
			if index+1 <= totalArgs && !strings.HasPrefix(parsedArgs[index+1], "-") {
				value = parsedArgs[index+1]
			}

			// Storing results from above in previously initialized map
			parsedArgMap[currentArg] = value
		}
	}

	// Merge Pluralith arg map with parsed arg map
	for arg, value := range pluralithArgs {
		parsedArgMap[arg] = value
	}

	// Create new arg array
	var finalArgs []string
	for arg, value := range parsedArgMap {
		if value == "" {
			finalArgs = append(finalArgs, "-"+arg)
		} else if strings.Contains(value, " ") {
			finalArgs = append(finalArgs, "-"+arg+"=\""+value+"\"")
		} else {
			finalArgs = append(finalArgs, "-"+arg+"="+value+"")
		}
	}

	// Returning slice of all args aswell as arg map
	return finalArgs
}
