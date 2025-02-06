package streamer

import (
	"bytes"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/hossein1376/queuer/config"
	"github.com/hossein1376/queuer/pkg/model"
	"github.com/hossein1376/queuer/pkg/order"
)

type streamer struct {
	client  *http.Client
	address string
	seed    int
}

// Run will generate random Orders and stream them to the worker.
func Run(cfg config.Streamer) {
	c := streamer{
		client:  &http.Client{Timeout: 10 * time.Second},
		address: cfg.Address,
		seed:    cfg.ProcessTimeSeed,
	}

	t := time.NewTicker(cfg.Interval)
	defer t.Stop()

	for range t.C {
		if err := c.send(); err != nil {
			slog.Error("send order", slog.Any("error", err))
		}
	}
}

func (s streamer) send() error {
	// generate a new order and decode it into bytes
	data := order.Encode(s.generateOrder())

	resp, err := s.client.Post(
		s.address, "application/json", bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	switch statusCode := resp.StatusCode; statusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return nil
	default:
		return fmt.Errorf("unexpected status code: %d", statusCode)

	}
}

func (s streamer) generateOrder() model.Order {
	return model.Order{
		OrderID:        rand.Int64(),
		Priority:       model.Priority(rand.IntN(2)),
		ProcessingTime: time.Duration(rand.IntN(s.seed)) * time.Second,
	}
}
