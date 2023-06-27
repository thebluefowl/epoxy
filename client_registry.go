package epoxy

import "sync"

var clientRegistry *ClientRegistry

type ClientRegistry struct {
	m sync.Map
}

func init() {
	clientRegistry = &ClientRegistry{
		m: sync.Map{},
	}
}

// Register registers a ReadWriter with the registry.  This is used to
// store the connection to the client.
func (r *ClientRegistry) Register(name string, value ReadWriter) {
	r.m.Store(name, value)
}

// Deregister removes a ReadWriter from the registry.  This is used to
// remove the connection to the client.  This also closes the underlying
// connection.

func (r *ClientRegistry) Deregister(name string) {
	v, ok := r.m.Load(name)
	if !ok {
		return
	}
	v.(ReadWriter).Close()
	r.m.Delete(name)
}

// Get returns a ReadWriter from the registry.  This is used to retrieve
// the connection to the client.
func (r *ClientRegistry) Get(name string) ReadWriter {
	value, ok := r.m.Load(name)
	if !ok {
		return nil
	}
	return value.(ReadWriter)
}

// Listen starts a goroutine that listens for messages from the client.
// This is used to process messages from the client.
func (r *ClientRegistry) Listen(name string) {
	value, ok := r.m.Load(name)
	if !ok {
		return
	}
	reader := value.(Reader)
	go func(r Reader) {
		for {
			b, err := r.Read()
			if err != nil {
				break
			}
			ProcessReply(b)
		}
	}(reader)
}
