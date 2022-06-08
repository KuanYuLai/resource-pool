package queue

import (
	"testing"
)

func TestPopAndPushBack(t *testing.T) {
	newQueue := NewQueue[int]()
	valueList := []int{1, 2, 3}

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
