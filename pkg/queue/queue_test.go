package queue

import (
	"testing"
	"time"
)

type T int

func TestPopAndPushBack(t *testing.T) {
	newQueue := NewQueue[T](time.Second)
	valueList := []T{1, 2, 3}

	// Test PushBack
	for i := range valueList {
		newQueue.PushBack(valueList[i])
	}
	if newQueue.Length() != 3 {
		t.Error("queue length should be 3")
	}

	// Test Pop
	for i := range valueList {
		v := newQueue.Pop()
		if v != valueList[i] {
			t.Errorf("value should be %d but got %d\n", valueList[i], v)
		}
	}
}

func TestIdleTimer(t *testing.T) {
	node1 := newNodePtr[T]()
	node2 := newNodePtr[T]()
	node3 := newNodePtr[T]()

	// link nodes
	node1.nextNode = node2
	node2.nextNode = node3
	node3.previousNode = node2
	node2.previousNode = node1

	newQueue := NewQueue[T](time.Second)
	newQueue.head = node1
	newQueue.tail = node3
	newQueue.length = 3

	acquired := make(chan struct{})

	// test time out
	idleTimer(acquired, 1*time.Second, newQueue, node2)
	if newQueue.length != 2 {
		t.Errorf("length should be 2 but got %d", newQueue.length)
	}

	if node1.nextNode != node3 || node3.previousNode != node1 {
		t.Errorf("nodes don't have proper relation")
	}

	printAllNodes(newQueue)

	// test acquired
	go idleTimer(acquired, 2*time.Second, newQueue, node3)
	time.Sleep(500 * time.Millisecond)
	acquired <- struct{}{}
	close(acquired)

	printAllNodes(newQueue)
	if newQueue.length != 2 {
		t.Errorf("length should be 2 but got %d", newQueue.length)
	}
}

func TestDeleteNode(t *testing.T) {
	node1 := newNodePtr[T]()
	node2 := newNodePtr[T]()
	node3 := newNodePtr[T]()

	// link nodes
	node1.nextNode = node2
	node2.nextNode = node3
	node3.previousNode = node2
	node2.previousNode = node1

	node1.value = 1
	node2.value = 2
	node3.value = 3

	newQueue := NewQueue[T](time.Second)
	newQueue.head = node1
	newQueue.tail = node3
	newQueue.length = 3
	newQueue.deleteNode(node2)

	if node1.nextNode != node3 || node3.previousNode != node1 {
		t.Errorf("nodes don't have proper relation")
	}
}
