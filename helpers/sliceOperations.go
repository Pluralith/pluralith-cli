package helpers

// - - - Collection of functions for common slice operations - - -

// Function to check if a given input appears in a slice
func ElementInSlice(element string, slice []string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}

// Function to get index of a given element in slice
func IndexInSlice(element string, slice []string) int {
	for key, value := range slice {
		if value == element {
			return key
		}
	}
	return -1
}

// Function to remove element from slice
func RemoveFromSlice(element string, slice []string) []string {
	index := IndexInSlice(element, slice)
	return append(slice[:index], slice[index+1:]...)
}
