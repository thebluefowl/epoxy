package registry

import (
	"sync"

	"github.com/thebluefowl/epoxy/internal/random"
	"github.com/thebluefowl/epoxy/internal/tunnel"
)

var Singleton *Registry

type Registry struct {
	m sync.Map
}

func NewRegistry() *Registry {
	return &Registry{
		m: sync.Map{},
	}
}

func (r *Registry) Get(name string) tunnel.Connection {
	value, ok := r.m.Load(name)
	if !ok {
		return nil
	}
	return value.(tunnel.Connection)
}

func (r *Registry) Register(c tunnel.Connection) error {
	id, err := random.URLSafeString(6)
	if err != nil {
		return err
	}
	// id = "localhost:8137"
	r.m.Store(id, c)
	c.Write([]byte(EPOXY_CTRL_ID + "=" + id))
	return nil
}

func (r *Registry) Deregister(id string) {
	value, ok := r.m.Load(id)
	if !ok {
		return
	}
	value.(tunnel.Connection).Close()
	r.m.Delete(id)
}
