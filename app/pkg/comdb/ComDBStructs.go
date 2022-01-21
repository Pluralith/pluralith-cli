package comdb

type ComDB struct {
	Events []ComDBEvent
}

type ComDBEvent struct {
	Receiver  string
	Timestamp int64
	Command   string
	Type      string
	Address   string
	Error     string `json:"Error,omitempty"`
	Message   string
	Instances []interface{}
	Path      string
	Received  bool
	Providers []string `json:"Providers,omitempty"`
}
