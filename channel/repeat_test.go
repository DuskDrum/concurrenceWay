package channel

import (
	"fmt"
	"math/rand"
	"testing"
)

// TestTakeSimpleCase  take 模式
func TestTakeSimpleCase(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	array := []any{1, 2, 3, 4}
	for num := range take(done, repeat(done, array...), 10) {
		fmt.Printf("%v ", num)
	}
}

// TestRepeatFnSimpleCase  take 模式, 重复调用函数
func TestRepeatFnSimpleCase(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	r := func() any {
		return rand.Int()
	}

	for num := range take(done, repeatFn(done, r), 10) {
		fmt.Printf("%v ", num)
	}
}

// TestToStringSimpleCase  take 模式, 重复调用函数
func TestToStringSimpleCase(t *testing.T) {

	done := make(chan interface{})
	defer close(done)
	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...", message)
}

func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for range toString(done, take(done, repeat(done, "a"), b.N)) {
	}
}

func BenchmarkTyped(b *testing.B) {
	repeat := func(done <-chan interface{}, values ...string) <-chan string {
		valueStream := make(chan string)
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}
	take := func(done <-chan interface{}, valueStream <-chan string, num int) <-chan string {
		takeStream := make(chan string)
		go func() {
			defer close(takeStream)
			for i := num; i > 0 || i == -1; {
				if i != -1 {
					i--
				}
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for range take(done, repeat(done, "a"), b.N) {
	}
}
