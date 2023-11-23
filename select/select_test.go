package _select

import (
	"fmt"
	"testing"
	"time"
)

// TestSimpleCase 简单的select--case。在go中for-select无限循环很常见
func TestSimpleCase(t *testing.T) {
	done := make(chan interface{})
	stringStream := make(chan string)

	go func() {
		defer close(stringStream)
		for _, s := range []string{"a", "b", "c"} {
			select {
			case <-done:
				fmt.Println("收到了done请求... ")
				return
			case stringStream <- s:
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case s := <-stringStream:
				if s != "" {
					fmt.Printf("%v ", s)
				}
			}
		}
	}()

	time.Sleep(5 * time.Second)

	close(done)

}
