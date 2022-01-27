package ux

import (
	"github.com/fatih/color"
)

// - - - Collection of UX related functions - - -

// Function to print Pluralith head art
func PrintHead() {
	color.HiBlue(` _
|_)|    _ _ |._|_|_ 
|  ||_|| (_||| | | |

`)
}

// Function to print formatted output
func PrintFormatted(text string, styling []string) {
	// Creating map with available styles
	styleMap := map[string]color.Attribute{
		"white": color.FgHiWhite,
		"blue":  color.FgHiBlue,
		"green": color.FgHiGreen,
		"red":   color.FgHiRed,
		"bold":  color.Bold,
	}

	// Defining new blank color object
	style := color.New()

	// Applying all styles found in "styling" arg
	for _, item := range styling {
		style.Add(styleMap[item])
	}

	// Printing passed text with perviously configured styles
	style.Printf(text)
}
