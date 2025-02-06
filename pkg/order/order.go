package order

import (
	"container/heap"
	"errors"
	"sync"

	"github.com/hossein1376/queuer/pkg/model"
)

var (
	ErrEmptyQueue = errors.New("queue is empty")
)

// order represents the queued [model.Order] instance.
type order struct {
	*model.Order

	// used and maintained by the heap tree
	index int
}

// Queue is a wrapper around [priorityQueue], making it safe for
// concurrent use, and providing type safety for the [Push] and [Pop]
// methods. It exposes only [Len], [Push] and [Pop] methods.
type Queue struct {
	pq *priorityQueue
	mu sync.Mutex
}

// NewQueue returns a new instance of [*Queue] which provides a
// concurrent and type safe priority queue.
func NewQueue() *Queue {
	pq := &priorityQueue{}
	heap.Init(pq)

	return &Queue{pq: pq}
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.pq.Len()
}

func (q *Queue) Push(o *model.Order) {
	q.mu.Lock()
	defer q.mu.Unlock()

	heap.Push(q.pq, &order{Order: o})
}

// Pop returns the first [model.Order] in the queue, or an error if the
// queue is empty.
func (q *Queue) Pop() (*model.Order, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.pq.Len() == 0 {
		return nil, ErrEmptyQueue
	}

	o := heap.Pop(q.pq).(*order)

	return o.Order, nil
}
