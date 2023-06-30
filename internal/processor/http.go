package processor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/thebluefowl/epoxy/internal/registry"
)

const (
	ExpoxyHeaderID = "X-Epoxy-ID"
)

type HTTP struct {
	requests sync.Map
	registry *registry.Registry
}

func NewHTTPProcessor(registry *registry.Registry) Processor {
	return &HTTP{
		requests: sync.Map{},
		registry: registry,
	}
}

func (p *HTTP) ProcessRequest(r Request) error {
	content, ok := r.Content.(*http.Request)
	if !ok {
		log.Print("invalid request type for HTTP processor: ", r.Content)
		return errors.New("invalid request type for HTTP processor")
	}
	content.Header.Set(ExpoxyHeaderID, r.ID)
	p.requests.Store(r.ID, r)

	b, err := httputil.DumpRequest(content, true)
	if err != nil {
		log.Print("error dumping request: ", err.Error())
		return err
	}

	writer := p.registry.Get(content.Host)
	if writer == nil {
		log.Println("no connection found")
		return errors.New("no connection found")
	}

	if err := writer.Write(b); err != nil {
		return err
	}
	return nil
}

func (p *HTTP) ProcessResponse(b []byte) error {
	res, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(b)), nil)
	if err != nil {
		return err
	}
	id := res.Header.Get(ExpoxyHeaderID)
	if id == "" {
		return errors.New("no Expoxy ID found in response")
	}

	r, ok := p.requests.Load(id)
	if !ok {
		return errors.New("no request found for response")
	}
	r.(Request).ResponseChan <- res
	return nil
}

type HTTPOpts struct {
	Port    int
	Timeout time.Duration
}

func (p *HTTP) Start(opts interface{}) error {
	o, ok := opts.(*HTTPOpts)
	if !ok {
		return errors.New("invalid options type for HTTP processor")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		request := NewRequest(r)
		if err := p.ProcessRequest(request); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		select {
		case response := <-request.ResponseChan:
			res := response.(*http.Response)
			for k, v := range res.Header {
				w.Header().Set(k, v[0])
			}
			w.WriteHeader(res.StatusCode)
			_, err := io.Copy(w, res.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case <-time.After(o.Timeout):
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}
	})
	log.Printf("Starting HTTP server on port %d", o.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", o.Port), mux)
}
