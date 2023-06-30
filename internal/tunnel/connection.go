package tunnel

type Connection interface {
	Reader
	Writer
}

type Reader interface {
	Close() error
	ReadAll() ([]byte, error)
}

type Writer interface {
	Close() error
	Write([]byte) error
	End() error
}
