package cond

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestSimpleTest cond的简单使用，cond用来实现等待、通知场景的并发问题。先newCond，然后wait，最后Single或者Broadcast
// 该程序成功地将所有10个项目添加到队列中(并且在它有机会将 前两项删除之前退出)
func TestSimpleTest(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}
	for i := 0; i < 10; i++ {
		// 上锁
		c.L.Lock()
		// 这里为什么是for
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}
