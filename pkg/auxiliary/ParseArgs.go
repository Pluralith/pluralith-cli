package auxiliary

import (
	"strings"
)

// - - - Code to manually parse flags - - -
// (Need to do it manually because cobra won't allow for unknown flags to pass through to terraform commands)

// Function to parse flags
func ParseArgs(args []string, pluralithArgs []string) ([]string, map[string]string) {
	// Instantiating clean slice to collect args
	var parsedArgs []string
	var cleanArgs []string
	parsedArgMap := make(map[string]string)

	// Cleaning up potential '=' in args
	for _, arg := range args {
		splitArg := strings.Split(arg, "=")
		parsedArgs = append(parsedArgs, splitArg...)
	}

	// Getting cleanArgs length to use in loop
	argLength := len(parsedArgs)
	// Creating arg map
	for index, arg := range parsedArgs {
		// Determining if current value is an arg or a value by checking for "-"
		if strings.HasPrefix(arg, "-") {
			// Initializing a few variables
			var nextArg string
			var value string

			// Clearing "-" prefix from args
			currentArg := arg[1:]

			// Checking if this is the last argument of the slice, to avoid error
			if argLength >= index+2 {
				// If there is a value beyond the current one -> Assign it to next arg
				nextArg = parsedArgs[index+1]
			} else {
				// Otherwise assign "true" to mark the current argument as present
				nextArg = "true"
			}

			// Checking if the next value (if given) contains "-" (which would make it an arg)
			if !strings.HasPrefix(nextArg, "-") {
				// If it doesn't it can be seen as value for current arg
				value = nextArg
			} else {
				// If it does, our current arg does not have a dedicated value and is a simple boolean arg
				value = "true"
			}

			// Checking if current arg is among pluralith args, removing it from args to be passed to Terraform
			if !ElementInSlice(arg, pluralithArgs) {
				if value == "true" {
					cleanArgs = append(cleanArgs, arg)
				} else {
					cleanArgs = append(cleanArgs, arg, nextArg)
				}
			}

			// Storing results from above in previously initialized map
			parsedArgMap[currentArg] = value
		}
	}

	// Returning slice of all args aswell as arg map
	return cleanArgs, parsedArgMap
}
