package auxiliary

type Filters struct {
	Replacement string
	Config      interface{}
	TempTargets []string
}

func (F *Filters) InitializeFilters() error {
	// functionName := "InitializeFilters"
	F.Replacement = "gatewatch"
	F.TempTargets = []string{"owner_id"}

	return nil
}

var FilterInstance = &Filters{}
