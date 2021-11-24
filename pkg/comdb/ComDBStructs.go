package comdb

type ComDB struct {
	Locked bool
	Events []interface{}
	Errors []interface{}
}

type Update struct {
	Receiver   string
	Timestamp  int64
	Command    string
	Event      string
	Address    string
	Attributes map[string]interface{}
	Path       string
	Received   bool
}
