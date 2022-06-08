package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/KuanYuLai/resource-pool_Dcard/pkg/pool"
)

func main() {
	type Company struct {
		name string
	}

	// create new pool
	fmt.Println("Step1: ")
	fmt.Println(`Creating new pool with settings:
	creator: return a Company struct pointer,
	maxIdleSize: 3
	maxIdleTime: 3 seconds`)

	p := pool.New(func(context.Context) (*Company, error) {
		return new(Company), nil
	}, 3, time.Second)
	time.Sleep(2 * time.Second)

	// acquire 4 items from the pool
	fmt.Println("\nStep2: ")
	fmt.Println("Acquire 4 items and assign company names 'Dcard', 'Otto', 'Facebook', and 'Twitter' respectively")
	item1, err := p.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	item2, err := p.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	item3, err := p.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	item4, err := p.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}

	// assign values to every item
	item1.name = "Dcard"
	item2.name = "Otto"
	item3.name = "Facebook"
	item4.name = "Twitter"

	fmt.Printf("item1: %+v\n", item1)
	fmt.Printf("item2: %+v\n", item2)
	fmt.Printf("item3: %+v\n", item3)
	fmt.Printf("item4: %+v\n", item4)
	time.Sleep(2 * time.Second)

	// release all items back to the pool
	fmt.Println("\nStep3: ")
	fmt.Println("Release all items back to the pool, but the pool should only contain 3 items")
	p.Release(item1)
	p.Release(item2)
	p.Release(item3)
	p.Release(item4)
	fmt.Printf("Number of idle items in the pool: %d\n", p.NumIdle())
	time.Sleep(2 * time.Second)

	fmt.Println("\nStep4: ")
	fmt.Println("Acquire one item from the pool, the item should contains the value 'Dcard' and the pool should only have 2 items left")
	item, err := p.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Acquired item: %+v\n", item)
	fmt.Printf("Current idle items in the pool: %d\n", p.NumIdle())
	time.Sleep(2 * time.Second)

	fmt.Println("\nStep5: ")
	fmt.Println("Wait 5 second to let all items in the pool expired, acquired a new item and  see if the item contains any value")
	for i := 1; i <= 5; i++ {
		fmt.Printf("%v...", i)
		time.Sleep(time.Second)
	}
	item, err = p.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("\nAcquired item: %+v\n", item)
	fmt.Printf("Current idle items in the pool: %d\n", p.NumIdle())
}
