package epoxy

import (
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func RecieverHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target := r.Header.Get(EpoxyHeaderTarget)
		if target == "" {
			http.Error(w, "Missing X-Target header", http.StatusBadRequest)
			return
		}

		writer := clientRegistry.Get(target)
		if writer == nil {
			http.Error(w, "Target not found", http.StatusNotFound)
			return
		}

		requestID := uuid.NewString()

		r.Header.Set(EpoxyHeaderRequestID, requestID)

		if err := writer.Write(r); err != nil {
			http.Error(w, "Failed to write request", http.StatusInternalServerError)
			return
		}

		// Register the request as inflight.
		inflightRegistry.Register(requestID, &InflightRequest{
			ResponseChan: make(chan *http.Response, 1),
		})

		select {
		case response := <-inflightRegistry.Get(requestID).ResponseChan:
			for k, v := range response.Header {
				w.Header().Set(k, v[0])
			}

			w.WriteHeader(response.StatusCode)

			_, err := io.Copy(w, response.Body)
			if err != nil {
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
				return
			}
		case <-time.After(time.Second * 5):
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

	}
}
