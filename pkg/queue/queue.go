package queue

import (
	"fmt"
	"sync"
	"time"
)

type (
	node[T any] struct {
		lock         *sync.Mutex
		previousNode *node[T]
		nextNode     *node[T]
		value        T
		acquired     chan struct{}
	}

	Queue[T any] struct {
		lock        *sync.Mutex
		head        *node[T]
		tail        *node[T]
		maxIdleTime time.Duration
		length      int
	}
)

func newNodePtr[T any]() *node[T] {
	node := &node[T]{
		lock:     new(sync.Mutex),
		acquired: make(chan struct{}),
	}
	return node
}

// NewQueue ...
func NewQueue[T any](maxIdleTime time.Duration) *Queue[T] {
	return &Queue[T]{
		lock:        new(sync.Mutex),
		maxIdleTime: maxIdleTime,
	}
}

// Pop remove the first node in queue
func (q *Queue[T]) Pop() T {
	firstNode := newNodePtr[T]()
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

func idleTimer[T any](acquired <-chan struct{}, timeLimit time.Duration, q *Queue[T], node *node[T]) {
	select {
	case <-time.After(timeLimit):
		q.deleteNode(node)
		return
	case <-acquired:
		return
	}
}

func (q *Queue[T]) deleteNode(node *node[T]) {
	q.lock.Lock()
	node.lock.Lock()

	// link left node to the right node
	if node.previousNode != nil {
		node.previousNode.lock.Lock()
		node.previousNode.nextNode = node.nextNode
		node.previousNode.lock.Unlock()
	}

	// link right node to the left node
	if node.nextNode != nil {
		node.nextNode.lock.Lock()
		node.nextNode.previousNode = node.previousNode
		node.nextNode.lock.Unlock()
	}

	// check if the node is head or tail node
	if q.isHead(node) {
		q.head = node.nextNode
	}
	if q.isTail(node) {
		q.tail = node.previousNode
	}
	node.lock.Unlock()
	q.length--
	q.lock.Unlock()
}

func printAllNodes[T any](q *Queue[T]) {
	node := q.head
	hasNext := true
	for hasNext {
		fmt.Printf("%v\t", node.value)
		node = node.nextNode
		hasNext = node != nil
	}
	fmt.Printf("\n")
}

func (q *Queue[T]) isTail(node *node[T]) bool {
	return q.tail == node
}

func (q *Queue[T]) isHead(node *node[T]) bool {
	return q.head == node
}
