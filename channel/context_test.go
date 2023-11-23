package channel

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

// TestDoneChannelCase done通道模式可以实现安全的控制程序，但是无法带有错误的额外信息
func TestDoneChannelCase(t *testing.T) {
	var wg sync.WaitGroup
	done := make(chan any)
	defer close(done)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()
	// 在这里阻塞， 所以此类里的两个goroutine会一起执行，各去等待5秒，算到这里也是5秒
	// 整个方法耗时5秒
	wg.Wait()
}

// TestContextCase 使用上下文实现协程的控制
func TestContextCase(t *testing.T) {
	var wg sync.WaitGroup
	// 定义cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreetingC(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewellC(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()
	wg.Wait()
}
