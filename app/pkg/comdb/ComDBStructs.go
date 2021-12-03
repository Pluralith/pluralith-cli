package comdb

type ComDB struct {
	Events []Event
	Errors []map[string]interface{}
}

type Event struct {
	Receiver  string
	Timestamp int64
	Command   string
	Type      string
	Address   string
	Instances []interface{}
	Path      string
	Received  bool
}
