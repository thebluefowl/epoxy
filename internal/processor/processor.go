package processor

type Processor interface {
	ProcessRequest(Request) error
	ProcessResponse([]byte) error
	Start(interface{}) error
}
