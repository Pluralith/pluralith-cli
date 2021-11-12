package communication

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
