package proxy

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
)

const (
	ExpoxyHeaderID = "X-Epoxy-ID"
)

type HTTP struct {
	Address string
}

func NewHTTPProxy(address string) Proxier {
	return &HTTP{
		Address: address,
	}
}

func (h *HTTP) Do(content []byte) ([]byte, error) {
	in, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(content)))
	if err != nil {
		return nil, err
	}
	id := in.Header.Get(ExpoxyHeaderID)

	out, err := http.NewRequest(
		in.Method,
		h.Address+in.URL.Path,
		in.Body,
	)

	if err != nil {
		return nil, err
	}

	for k, v := range in.Header {
		out.Header[k] = v
	}

	res, err := http.DefaultClient.Do(out)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	res.Header.Set(ExpoxyHeaderID, id)

	return httputil.DumpResponse(res, true)
}
