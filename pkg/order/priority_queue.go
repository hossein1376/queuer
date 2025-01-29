package order

// priorityQueue represents a queue of [*order]s, via implementing the
// [heap.Interface].
//
// Orders with higher [model.Priority] are proceed before the ones with
// lower [model.Priority], thus making it a priority queue with a heap
// tree under the hood.
type priorityQueue []*order

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push accepts an instance of [*order], and add it to the queue.
func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*order)
	item.index = n
	*pq = append(*pq, item)
}

// Pop returns the first item in queue. Its type will be of [*order].
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // // let gc clean it up
	item.index = -1 // just in case
	*pq = old[0 : n-1]

	return item
}
