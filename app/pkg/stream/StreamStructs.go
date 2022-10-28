package stream

type DecodedEvent struct {
	Command    string
	Type       string
	ParsedType string
	Address    string
	Message    string
	Outputs    map[string]interface{}
}
