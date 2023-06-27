package epoxy

import "net/http"

type ReadWriter interface {
	Reader
	Writer
}

type Reader interface {
	// Read() ([]byte, error)
	Close()
	Read() (*http.Response, error)
}

type Writer interface {
	// Write([]byte) error
	Close()
	Write(*http.Request) error
}
