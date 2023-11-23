package pipeline

import (
	"fmt"
	"testing"
)

func TestPipelineSimpleCase(t *testing.T) {
	// pipeline有批处理和流处理之分，每次返回一整个切片，也可以每次传入一个值，返回一个值。比较推荐这样处理，这样更加灵活可拓展(扇出fan-out)
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	// 所有记录先乘以2，再+1
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}

	// 所有记录先乘以2，再+1,再乘以3
	for _, v := range multiply(add(multiply(ints, 2), 1), 3) {
		fmt.Println(v)
	}
}

// pipeline非常适合配合channel使用
func TestPipelineChannelCase(t *testing.T) {
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
