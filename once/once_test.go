package once

import (
	"fmt"
	"sync"
	"testing"
)

// TestSimpleCase 无论调用多少次once.Do，最后实际上只会执行一次，输出结果为1
func TestSimpleCase(t *testing.T) {
	var count int
	increment := func() {
		count++
	}
	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}
	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

// TestDuplicateDoCase 同一个once，调用多个不同的Do最终也只会有第一个生效，如果想要两个同时生效需要定义两个once
func TestDuplicateDoCase(t *testing.T) {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }
	var once sync.Once
	once.Do(increment)
	once.Do(decrement)
	fmt.Printf("Count :%d\n", count)
}

// TestDuplicateDoResolverCase 一加一减，最后才是0
func TestDuplicateDoResolverCase(t *testing.T) {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }
	var once sync.Once
	var once2 sync.Once
	once.Do(increment)
	once2.Do(decrement)
	fmt.Printf("Count :%d\n", count)
}

// TestCycleDoCase 循环调用导致死锁问题， 会报错：fatal error: all goroutines are asleep - deadlock!
func TestCycleDoCase(t *testing.T) {
	var onceA, onceB sync.Once
	var initB func()
	initA := func() { onceB.Do(initB) }
	initB = func() { onceA.Do(initA) }
	onceA.Do(initA)
}
