package channel

import (
	"fmt"
	"testing"
)

// TestOrChannelSimpleCase  or-channel 模式
// 想把一个channel 一变二；以便 将它们发送到 代码的两个不同独立区域中。
func TestTeeChannelSimpleCase(t *testing.T) {
	done := make(chan any)
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2, 3, 4), 4))
	// 阻塞的原因居然在这里，for循环out1 会导致out2一直不被消费
	for val1 := range out1 {
		fmt.Printf("out1: %v , out2: %v \n ", val1, <-out2)
	}
}
