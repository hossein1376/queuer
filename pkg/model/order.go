package model

import (
	"time"
)

// Order defines structure of an order.
type Order struct {
	OrderID        int64
	ProcessingTime time.Duration
	Priority       Priority
	Status         Status
}
