package communication

type CommunicationDB struct {
	Locked bool
	Events []interface{}
	Errors []interface{}
}
