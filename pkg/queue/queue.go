package queue

import (
	"sync"
)

type (
	node[T any] struct {
		lock         *sync.Mutex
		previousNode *node[T]
		nextNode     *node[T]
		value        T
	}

	Queue[T any] struct {
		lock   *sync.Mutex
		head   *node[T]
		tail   *node[T]
		length int
	}
)

// NewQueue ...
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		lock: new(sync.Mutex),
	}
}

// Pop remove the first node in queue
func (q *Queue[T]) Pop() T {
	firstNode := &node[T]{
		lock: new(sync.Mutex),
	}
	q.lock.Lock()
	firstNode = q.head
	q.head = firstNode.nextNode

	// unlink firstNode to chain
	if firstNode.nextNode != nil {
		firstNode.lock.Lock()
		firstNode.nextNode.lock.Lock()
		firstNode.nextNode.previousNode = nil
		firstNode.nextNode.lock.Unlock()
		firstNode.lock.Unlock()
	}

	q.length--
	q.lock.Unlock()
	return firstNode.value
}

// Pop remove the first node in queue
func (q *Queue[T]) PushBack(value T) {
	q.lock.Lock()
	lastNode := &node[T]{
		lock:         new(sync.Mutex),
		previousNode: q.tail,
		value:        value,
	}
	if q.length == 0 {
		q.head = lastNode
		q.tail = lastNode
	} else {
		q.tail.nextNode = lastNode
	}
	q.tail = lastNode
	q.length++
	q.lock.Unlock()
}

func (q *Queue[T]) Length() int {
	q.lock.Lock()
	l := q.length
	q.lock.Unlock()
	return l
}

func (q *Queue[T]) DeleteNode(node *node[T]) {
}

func (q *Queue[T]) isTail(node *node[T]) bool {
	return q.tail == node
}

func (q *Queue[T]) isHead(node *node[T]) bool {
	return q.head == node
}
