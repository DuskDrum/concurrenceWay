package pipeline

import (
	"fmt"
	"sync"
	"testing"
)

// 扇出(Fan-out)是一个术语，用于描述启动多个goroutines以处理来自管道的输入的过程，并且扇入 (fan-in)是描述将多个结果组合到一个通道中的过程的术语。
func TestPipelineFallInCase(t *testing.T) {
	fallIn := func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})
		multiplex := func(c <-chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}
		//从所有 的 channel 里取值
		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}
		// 等待所有的读操作结束
		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()
		return multiplexedStream
	}

	fallIn(nil, nil)

	// done channel用来防止goroutine 泄露
	// generator方法，用来将数组、切片转化为channel。离散值转换为 channel 上的值流
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				// select 配合done也是为了避免泄露goroutine
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}
	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addStream := make(chan int)
		go func() {
			defer close(addStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addStream <- i + additive:
				}
			}
		}()
		return addStream
	}
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)

	// 先*2，再加1，再乘2
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)
	// 所有记录先乘以2，再+1,再乘以3
	for v := range pipeline {
		fmt.Println(v)
	}
}
