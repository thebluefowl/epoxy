package epoxy

import (
	"fmt"
	"net/http"
	"sync"
)

var inflightRegistry *InflightRegistry

// InflightRequest is a request that is waiting for a response from
// the underlying communication channel.
type InflightRequest struct {
	ID           string
	ResponseChan chan *http.Response
}

// InflightRegistry is a registry of requests that are waiting for
// a response from the underlying communication channel.  This is
// used to correlate the response with the original request.
type InflightRegistry struct {
	m sync.Map
}

// Registers a new InflightRequest with the registry.
func (r *InflightRegistry) Register(id string, request *InflightRequest) {
	r.m.Store(id, request)
}

// Deletes an InflightRequest from the registry.
func (r *InflightRegistry) Deregister(id string) {
	r.m.Delete(id)
}

// Get returns an InflightRequest from the registry.
func (r *InflightRegistry) Get(id string) *InflightRequest {
	value, ok := r.m.Load(id)
	if !ok {
		return nil
	}
	return value.(*InflightRequest)
}

// init initializes the inflightRegistry singleton.
func init() {
	fmt.Println("initialized inflightRegistry")
	inflightRegistry = &InflightRegistry{m: sync.Map{}}
}
