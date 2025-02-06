package stats

import (
	"fmt"
	"sync"
	"time"

	"github.com/hossein1376/queuer/pkg/model"
)

// Stats stores the worker statics of the processed orders.
type Stats struct {
	counter map[model.Status]uint64
	mu      sync.Mutex
}

func New() *Stats {
	return &Stats{counter: make(map[model.Status]uint64)}
}

func (s *Stats) Add(status model.Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter[status] += 1
}

func (s *Stats) Print(pending int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf(
		"@%s\n\tPending: %d\n\tProcessed: %d\n\tFailed: %d\n",
		time.Now().Format(time.RFC3339),
		pending,
		s.counter[model.Processed],
		s.counter[model.Failed],
	)
}
