package worker

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hossein1376/querier/pkg/order"
)

// serve starts an HTTP server on the provided address.
func (p *pool) serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", p.handler)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv.ListenAndServe()
}

// handler receives a new order, properly handles reading the request's
// body and decoding the order data, and will subsequently queue it.
func (p *pool) handler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		writeResponse(
			w,
			http.StatusBadRequest,
			fmt.Errorf("read request: %w", err),
		)
		return
	}

	ord, err := order.Decode(b)
	if err != nil {
		writeResponse(
			w,
			http.StatusBadRequest,
			fmt.Errorf("decode order: %w", err),
		)
		return
	}

	// Order's status will be its zero value, [model.Pending], so no
	// subsequent checks are required.

	// if orders channel is full, this operation will be blocking. To
	// avoid that, create a separate goroutine for queuing each order.
	go p.queueOrder(ord)

	writeResponse(w, http.StatusOK, nil)
}

func writeResponse(w http.ResponseWriter, code int, msg error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	if msg != nil {
		w.Write([]byte(msg.Error()))
	}
}
