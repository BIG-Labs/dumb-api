package utils

import (
	"math/big"
)

// Queue represents a FIFO data structure
type Queue[T any] struct {
	elements []T
}

// NewQueue creates a new Queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{elements: make([]T, 0)}
}

// Enqueue adds an element to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	q.elements = append(q.elements, value)
}

// Dequeue removes an element from the front of the queue and returns it
func (q *Queue[T]) Dequeue() T {
	var result T
	if len(q.elements) == 0 {
		return result
	}
	result = q.elements[0]
	q.elements = q.elements[1:]
	return result
}

// IsEmpty returns true if the queue is empty (doesn't work)
func (q *Queue[T]) IsEmpty() bool {
	return len(q.elements) == 0
}

// Len returns the lenght of the queue
func (q *Queue[T]) Len() int {
	return len(q.elements)
}

type Item struct {
	ID       uint
	Distance *big.Float
	Index    int
}

// PriorityQueue represents an ordered data structure
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Distance.Cmp(pq[j].Distance) < 0
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}
