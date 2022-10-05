package ci

type ApplyEventCosts struct {
	Hourly  float64 `json:"hourly"`
	Monthly float64 `json:"monthly"`
}

type ApplyEvent struct {
	Level     string `json:"@level"`
	Message   string `json:"@message"`
	Module    string `json:"@module"`
	Timestamp string `json:"@timestamp"`
	Type      string `json:"type"`
	Hook      struct {
		Action         string `json:"action"`
		IDKey          string `json:"id_key"`
		IDValue        string `json:"id_value"`
		ElapsedSeconds int    `json:"elapsed_seconds"`
		Resource       struct {
			Addr            string          `json:"addr"`
			Module          string          `json:"module"`
			Resource        string          `json:"resource"`
			ImpliedProvider string          `json:"implied_provider"`
			ResourceType    string          `json:"resource_type"`
			ResourceName    string          `json:"resource_name"`
			ResourceKey     interface{}     `json:"resource_key"`
			Costs           ApplyEventCosts `json:"costs"`
		} `json:"resource"`
	} `json:"hook"`
}
