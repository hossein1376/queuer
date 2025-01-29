package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/hossein1376/querier/config"
	"github.com/hossein1376/querier/pkg/model"
)

// Run creates a new worker pool, runs the HTTP server, and starts
// processing orders.
func Run(cfg config.Worker) {
	p := newPool(cfg.MaxJobs, cfg.Timeout)
	go p.printStats()
	go p.start()

	if err := p.serve(cfg.Address); err != nil {
		// failure to start the server is a serious issue, worthy of
		// panicking.
		panic(fmt.Errorf("start worker server: %w", err))
	}
}

// processOrder simulates the act of order's processing, sleeping for
// the duration of ProcessingTime, or canceling the context after
// the threshold set by timeout is passed, whichever happens first.
func (p *pool) processOrder(ord *model.Order) model.Status {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	select {
	case <-time.After(ord.ProcessingTime):
		return model.Processed
	case <-ctx.Done():
		return model.Failed
	}
}

func (p *pool) printStats() {
	for {
		p.stats.Print(p.queue.Len())
		time.Sleep(2 * time.Second)
	}
}
