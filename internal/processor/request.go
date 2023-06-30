package processor

import "github.com/google/uuid"

type Request struct {
	ID           string
	ResponseChan chan interface{}
	Content      interface{}
}

func NewRequest(content interface{}) Request {
	return Request{
		ID:           uuid.NewString(),
		ResponseChan: make(chan interface{}),
		Content:      content,
	}
}
