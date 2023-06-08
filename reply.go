package main

import (
	"net/http"
)

func ProcessReply(res *http.Response) {
	requestID := res.Header.Get(HeaderRequestID)
	if requestID == "" {
		return
	}
	if inflightRegistry.Get(requestID) == nil {
		return
	}

	inflightRegistry.Get(requestID).ResponseChan <- res
}
