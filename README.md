# resource-pool

# 設計思路:
在設計管理 resource pool 中使用到的 queue 時, 原本想將 delete idle node 這個功能交由一個 goroutine 來管理，但我認為在計算每個 node 剩餘時間的功能實作上會相當複雜， 所以我選擇將每個 node 的 idle time 都由獨立的 goroutine 來計算，並且當時間到時，這個 goroutine 也會將它負責的 node 清除。在每個獨立的 goroutine 中，都有接受系統信號的 channel，這樣在使用者終止程式時，可以確保這些 goroutine 都會被即時關閉，而非等到時間到時 goroutine 才被關閉。因為需要刪除 idle node 的緣故，我選擇使用 linked list 而非使用 slice 來實作 queue，原因為 linked list 再刪除 node 的操作上非常簡單，也不會消耗太多效能。另外，這個 resource pool 不是一開始就將 pool 中能容納的 resource 建立起來，而是依照使用者使用數量慢慢增加，因此，如果使用 slice 實作，當要新增 slice 的空間時，需要的效能會隨著 slice 的大小越來越多。所以我選擇使用 linked list 來實作 queue。

在這個 resource pool 中, 我沒有將 maxIdleSize 的數量當做使用者可以拿取 resource 的最大數量。原因為，在設定 maxIdleSize 時，只能大略估算使用者需要的數量，無法非常精準，所以當需求超過 maxIdleSize 時，使用者依然可以要的到 resource，不會因為 resource 數量限制而卡住。這樣的作法也有可能發生給太多 resource 導致機器當機，但我認為如果遇到這種情況，也是一個系統給出的信號，有可能是使用者要了太多不必要的資源，也有可能是設定 maxIdleSize 時低估了使用者的需求。

# Usage
> :warning: This program require **go 1.18 or newer**!

Run demo function:
```sh
go run main.go
```
The demo contains **five steps**:
1. Create a new pool with the following settings:
    * creator: return a Compnay struct pointer
    * maxIdleSize: 3
    * maxIdleTime: 3 seconds
2. Acquire items from pool and assign value to each item
3. Release all items back to the pool
4. Acquire one item and check for the value
5. Acquire a item and check for it's value after all items exist in the pool expired
