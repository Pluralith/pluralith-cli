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

	argLength := len(parsedArgs) // Getting cleanArgs length to use in loop

	// Creating arg map
	for index, arg := range parsedArgs {
		// Determining if current value is an arg or a value by checking for "-"
		if strings.HasPrefix(arg, "-") {
			// Initializing a few variables
			var nextArg string
			var value string

			currentArg := arg[1:] // Clearing "-" prefix from args

			// Checking if this is the last argument of the slice, to avoid error
			if argLength >= index+2 {
				nextArg = parsedArgs[index+1] // If there is a value beyond the current one -> Assign it to next arg
			} else {
				nextArg = "true" // Otherwise assign "true" to mark the current argument as present
			}

			// Checking if the next value (if given) contains "-" (which would make it an arg)
			if !strings.HasPrefix(nextArg, "-") {
				value = nextArg // If it doesn't it can be seen as value for current arg
			} else {
				value = "true" // If it does, our current arg does not have a dedicated value and is a simple boolean arg
			}

			// Storing results from above in previously initialized map
			parsedArgMap[currentArg] = value
		}
	}

	// Merge pluralith arg map with parsed arg map
	for arg, value := range pluralithArgs {
		parsedArgMap[arg] = value
	}

	// Create new arg array
	var finalArgs []string
	for arg, value := range parsedArgMap {
		if strings.Contains(value, " ") {
			finalArgs = append(finalArgs, "-"+arg+"=\""+value+"\"")
		} else {
			finalArgs = append(finalArgs, "-"+arg+"="+value+"")
		}
	}

	// Returning slice of all args aswell as arg map
	return finalArgs
}
