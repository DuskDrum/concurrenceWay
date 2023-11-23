package goroutine

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

// 一个新创建的goroutine被赋予了几千字节，这在大部分情况都是足够的。当它不运行时，Go语言就会自动增长(缩小)存储堆栈的内存，允许多个goroutine存在适当的内存中。每个函数调用CPU的开销平均为3个廉价指令。在同一个地址空间中创建成千上万的goroutine是可行的。 如果goroutine只是线程，系统的资源消耗会更小。
// 创建了10000个goroutine，2.77kb。非常轻量
func TestSize(t *testing.T) {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}
	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }
	const numGoroutines = 1e4
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}

// BenchmarkContextSwitch 测试上下文切换的耗时，使用基准测试。测试OS线程和goroutine之间切换的上下文的相对性能
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})
	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}
	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}
	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin)
	wg.Wait()
}
