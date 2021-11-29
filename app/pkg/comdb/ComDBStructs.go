package comdb

type ComDB struct {
	Locked bool
	Events []Event
	Errors []map[string]interface{}
}

type Event struct {
	Receiver   string
	Timestamp  int64
	Command    string
	Type       string
	Address    string
	Attributes map[string]interface{}
	Path       string
	Received   bool
}
