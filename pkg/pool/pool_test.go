package pool

import (
	"context"
	"testing"
	"time"
)

type testStruct struct {
	name string
}

func TestBasicPoolFunction(t *testing.T) {
	newPool := New(func(ctx context.Context) (*testStruct, error) {
		v := new(testStruct)
		return v, nil
	}, 2, time.Second)

	item1, err := newPool.Acquire(context.Background())
	if err != nil {
		t.Error(err)
	}
	item1.name = "item1"

	item2, err := newPool.Acquire(context.Background())
	if err != nil {
		t.Error(err)
	}
	if item2.name != "" {
		t.Errorf("name value should be 'item1' but got %s", item2.name)
	}
	item2.name = "item2"
	newPool.Release(item1)
	newPool.Release(item2)
	if newPool.NumIdle() != 2 {
		t.Errorf("pool length should be 2 but got %d", newPool.NumIdle())
	}
	item, err := newPool.Acquire(context.Background())
	if item.name != "item1" {
		t.Errorf("name value should be 'item1' but got %s", item.name)
	}
	if newPool.NumIdle() != 1 {
		t.Errorf("pool length should be 1 but got %d", newPool.NumIdle())
	}
}

func TestMaxIdleSize(t *testing.T) {
	newPool := New(func(ctx context.Context) (*testStruct, error) {
		v := new(testStruct)
		return v, nil
	}, 40, time.Second)

	itemList := make([]*testStruct, 50)
	for i := 0; i < 50; i++ {
		item, err := newPool.Acquire(context.Background())
		if err != nil {
			t.Error(err)
		}
		itemList = append(itemList, item)
	}

	for i := 0; i < 50; i++ {
		newPool.Release(itemList[i])
	}

	if newPool.NumIdle() != 40 {
		t.Errorf("NumIdle should be 40 but got %d", newPool.NumIdle())
	}
}
