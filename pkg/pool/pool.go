package pool

import (
	"context"
	"sync"
	"time"

	"github.com/KuanYuLai/resource-pool_Dcard/pkg/queue"
)

type Pool[T any] interface {
	// This creates or returns a ready-to-use item from the resource pool
	Acquire(context.Context) (T, error)
	// This releases an active resource back to the resource pool
	Release(T)
	// This returns the number of idle items
	NumIdle() int
}

type ResourcePool[T any] struct {
	lock        *sync.Mutex
	resource    func(context.Context) (T, error)
	idlePool    *queue.Queue[T]
	maxIdleSize int
	maxIdleTime time.Duration
}

func (r *ResourcePool[T]) Acquire(ctx context.Context) (T, error) {
	r.lock.Lock()
	// create new item if there isn't any idle item in the pool
	if r.idlePool.Length() == 0 {
		r.lock.Unlock()
		return r.resource(ctx)
	}
	item := r.idlePool.Pop()
	r.lock.Unlock()
	return item, nil
}

func (r *ResourcePool[T]) Release(item T) {
	r.lock.Lock()
	if r.idlePool.Length() >= r.maxIdleSize {
		// clear variable for gc
		item = *new(T)
		r.lock.Unlock()
		return
	}
	r.idlePool.PushBack(item)
	r.lock.Unlock()
}

func (r *ResourcePool[T]) NumIdle() int {
	r.lock.Lock()
	l := r.idlePool.Length()
	r.lock.Unlock()
	return l
}

// New ..
// creator is a function called by the pool to create a resource.
// maxIdleSize is the number of maximum idle items kept in the pool
// maxIdleTime is the maximum idle time for an idle item to be swept from the pool
func New[T any](
	creator func(context.Context) (T, error),
	maxIdleSize int,
	maxIdleTime time.Duration,
) Pool[T] {
	return &ResourcePool[T]{
		lock:        new(sync.Mutex),
		resource:    creator,
		maxIdleSize: maxIdleSize,
		maxIdleTime: maxIdleTime,
		idlePool:    queue.NewQueue[T](),
	}
}
