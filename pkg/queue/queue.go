package queue

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	// node store value of T type
	node[T any] struct {
		lock         *sync.Mutex
		previousNode *node[T]
		nextNode     *node[T]
		value        T
		acquired     chan struct{}
	}

	// Queue a queue design for resource pool
	Queue[T any] struct {
		lock        *sync.Mutex
		head        *node[T]
		tail        *node[T]
		maxIdleTime time.Duration
		length      int
	}
)

// newNodePtr returns a new node pointer
func newNodePtr[T any]() *node[T] {
	node := &node[T]{
		lock:     new(sync.Mutex),
		acquired: make(chan struct{}),
	}
	return node
}

// NewQueue returns a new queue
func NewQueue[T any](maxIdleTime time.Duration) *Queue[T] {
	return &Queue[T]{
		lock:        new(sync.Mutex),
		maxIdleTime: maxIdleTime,
	}
}

// Pop remove the first node in the queue
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

	// stop the idleTimer to prevent node deletion
	firstNode.acquired <- struct{}{}
	close(firstNode.acquired)
	q.length--
	q.lock.Unlock()
	return firstNode.value
}

// PushBack adds a node to the end of the queue
func (q *Queue[T]) PushBack(value T) {
	q.lock.Lock()
	lastNode := newNodePtr[T]()
	lastNode.previousNode = q.tail
	lastNode.value = value

	if q.length == 0 {
		q.head = lastNode
		q.tail = lastNode
	} else {
		q.tail.nextNode = lastNode
	}

	// create a goroutine to time the node's idle time
	go idleTimer(q, lastNode)

	q.tail = lastNode
	q.length++
	q.lock.Unlock()
}

// Length return the length of the queue
func (q *Queue[T]) Length() int {
	q.lock.Lock()
	l := q.length
	q.lock.Unlock()
	return l
}

// idleTimer time the idle time from the queue's setting. Delete the given node
// when the time is up.
func idleTimer[T any](q *Queue[T], node *node[T]) {
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), q.maxIdleTime)
	defer close(osSignal)
	defer cancel()
	select {
	case <-ctx.Done():
		q.deleteNode(node)
		return
	case <-node.acquired:
		return
	case <-osSignal:
		return
	}
}

// deleteNode delete the given node from the queue
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

// printAllNodes prints the entire queue.
// func printAllNodes[T any](q *Queue[T]) {
// 	node := q.head
// 	hasNext := true
// 	for hasNext {
// 		fmt.Printf("%v\t", node.value)
// 		node = node.nextNode
// 		hasNext = node != nil
// 	}
// 	fmt.Printf("\n")
// }

// isTail check if the given node is the last item in the queue
func (q *Queue[T]) isTail(node *node[T]) bool {
	return q.tail == node
}

// isHead check if the given node is the first item in the queue
func (q *Queue[T]) isHead(node *node[T]) bool {
	return q.head == node
}
