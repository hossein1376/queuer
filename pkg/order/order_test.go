package order_test

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hossein1376/queuer/pkg/model"
	"github.com/hossein1376/queuer/pkg/order"
)

func TestQueue_Len(t *testing.T) {
	q := order.NewQueue()

	q.Push(newOrder(model.Normal))
	q.Push(newOrder(model.Normal))
	q.Pop()
	q.Push(newOrder(model.Normal))

	assert.Equal(t, q.Len(), 2)
}

func TestQueue_Pop(t *testing.T) {
	a := assert.New(t)
	q := order.NewQueue()

	lowerPriority := newOrder(model.Normal)
	q.Push(lowerPriority)

	higherPriority := newOrder(model.High)
	q.Push(higherPriority)

	first, err := q.Pop()
	a.NoError(err)
	a.Equal(higherPriority, first)

	second, err := q.Pop()
	a.NoError(err)
	a.Equal(lowerPriority, second)

	item, err := q.Pop()
	a.ErrorIs(err, order.ErrEmptyQueue)
	a.Nil(item)
}

func TestQueue_Push(t *testing.T) {
	a := assert.New(t)
	q := order.NewQueue()
	ord := newOrder(model.Normal)

	q.Push(ord)
	a.Equal(q.Len(), 1)

	got, err := q.Pop()
	a.NoError(err)
	a.Equal(ord, got)
}

func newOrder(p model.Priority) *model.Order {
	return &model.Order{
		Priority:       p,
		OrderID:        rand.Int64(),
		ProcessingTime: time.Duration(rand.IntN(11)) * time.Second,
		Status:         model.Status(rand.IntN(3)),
	}
}
