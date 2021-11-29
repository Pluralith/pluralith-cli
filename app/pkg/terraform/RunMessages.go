package terraform

var RunMessages = map[string]interface{}{
	"plan":    []string{" Running plan ⇢ Apply it in the Pluralith UI\n", "Waiting for Confirmation ⇢ Apply plan in the Pluralith UI", "Apply Confirmed", "Apply Canceled"},
	"apply":   []string{" Running apply ⇢ Confirm it in the Pluralith UI\n", "Waiting for Confirmation ⇢ Confirm apply in the Pluralith UI", "Apply Confirmed", "Apply Canceled"},
	"destroy": []string{" Running destroy ⇢ Confirm it in the Pluralith UI\n", "Waiting for Confirmation ⇢ Confirm destroy in the Pluralith UI", "Destroy Confirmed", "Destroy Canceled"},
}
