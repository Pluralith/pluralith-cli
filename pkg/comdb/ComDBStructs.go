package comdb

type ComDB struct {
	Locked bool
	Events []Update
	Errors []map[string]interface{}
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

// Receiver  string
// Timestamp  float64
// Command  string
// Event  string
// Address  string
// Attributes  map[string]interface {}
// Path  string
// Received  bool
