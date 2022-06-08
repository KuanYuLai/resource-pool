package pool

import (
	"context"
	"time"
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
	resource    func(context.Context) (T, error)
	idlePool    []T
	maxIdleSize int
	maxIdleTime time.Duration
}

func (r *ResourcePool[T]) Acquire(ctx context.Context) (T, error) {
	if len(r.idlePool) == 0 {
		return r.resource(ctx)
	}
	item := r.idlePool[0]
	r.idlePool = r.idlePool[1:]
	return item, nil
}

func (r *ResourcePool[T]) Release(item T) {
	if len(r.idlePool) < r.maxIdleSize {
		r.idlePool = append(r.idlePool, item)
	}
}

func (r *ResourcePool[T]) NumIdle() int {
	return len(r.idlePool)
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
		resource:    creator,
		maxIdleSize: maxIdleSize,
		maxIdleTime: maxIdleTime,
		idlePool:    make([]T, 0),
	}
}
