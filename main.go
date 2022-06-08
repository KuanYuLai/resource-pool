package main

import (
	"context"
	"fmt"
	"time"

	"github.com/KuanYuLai/resource-pool_Dcard/pkg/pool"
)

func main() {
	type testStruct struct {
		name string
	}

	newPool := pool.New(func(ctx context.Context) (*testStruct, error) {
		v := new(testStruct)
		return v, nil
	}, 1, time.Hour)

	fmt.Println(newPool.NumIdle())

	s1, _ := newPool.Acquire(context.Background())
	fmt.Println(s1)
	s1.name = "Used item"
	fmt.Println(s1)
	newPool.Release(s1)
	fmt.Println(newPool.NumIdle())
	s1, _ = newPool.Acquire(context.Background())
	fmt.Println(s1)
	fmt.Println(newPool.NumIdle())
}
