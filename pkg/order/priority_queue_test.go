package order

import (
	"container/heap"
)

var _ heap.Interface = &priorityQueue{}
