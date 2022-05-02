package terraform

var RunMessages = map[string]interface{}{
	"plan":    []string{" Running Plan ⇢ Apply in the Pluralith UI\n", "Waiting for Confirmation ⇢ Apply Plan in the Pluralith UI", "Apply Confirmed", "Apply Canceled"},
	"apply":   []string{" Running Apply ⇢ Confirm in the Pluralith UI\n", "Waiting for Confirmation ⇢ Confirm Apply in the Pluralith UI", "Apply Confirmed", "Apply Canceled"},
	"destroy": []string{" Running Destroy ⇢ Confirm in the Pluralith UI\n", "Waiting for Confirmation ⇢ Confirm Destroy in the Pluralith UI", "Destroy Confirmed", "Destroy Canceled"},
}
