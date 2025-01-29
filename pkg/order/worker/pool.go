package worker

import (
	"time"

	"github.com/hossein1376/querier/pkg/model"
	"github.com/hossein1376/querier/pkg/order"
	"github.com/hossein1376/querier/pkg/order/stats"
)

// pool controls the goroutines worker pool, run the order processing
// job, and stores the statics.
type pool struct {
	orders  chan struct{}     // pending requests queue
	jobs    chan struct{}     // registry of running workers
	result  chan model.Status // reporting the processing result
	timeout time.Duration     // order processing deadline
	stats   *stats.Stats      // stores the statics
	queue   *order.Queue      // the priority queue
}

// newPool creates a new instance of the worker pool. With zero or a
// negative jobsCount, no limitations will be placed on the number of
// worker goroutines.
func newPool(jobsCount int, timeout time.Duration) *pool {
	var jobs chan struct{}
	if jobsCount > 0 {
		jobs = make(chan struct{}, jobsCount)
	}

	return &pool{
		queue:   order.NewQueue(),
		stats:   stats.New(),
		timeout: timeout,
		jobs:    jobs,
		result:  make(chan model.Status),
		orders:  make(chan struct{}),
	}
}

// start the worker pool, create a new goroutine for each order, and
// update the status statics.
func (p *pool) start() {
	go func() {
		for status := range p.result {
			// receive order status and update the statics
			p.stats.Add(status)
		}
	}()

	for range p.orders {
		if p.jobs != nil {
			// insert a new job instance, or block until there's a spot
			p.jobs <- struct{}{}
		}
		go p.jobHandler()
	}

}

// queueOrder pushes the order into the priority queue, and sends a
// signal indicating that a new order was received.
func (p *pool) queueOrder(ord *model.Order) {
	p.queue.Push(ord)
	p.orders <- struct{}{}
}

// jobHandler pops the first item in the queue, and process it.
func (p *pool) jobHandler() {
	defer func() {
		if p.jobs != nil {
			// indicates that the job is finished
			<-p.jobs
		}
	}()

	// Pop is concurrent safe, and checks for the queue's length
	ord, err := p.queue.Pop()
	if err != nil {
		// Since there is a 1:1 ratio between incoming orders and the
		// number of running goroutines, this err is not likely to ever
		// happen. Still, like any good programmer, we hope for the best
		// and prepare for the worst.
		return
	}

	status := p.processOrder(ord)
	ord.Status = status

	// Here, we should probably do something with the order. But since
	// we care only about the number of processed orders, we pass down
	// its status and ignore the rest.
	p.result <- status
}
