package channel

import (
	"fmt"
	"testing"
	"time"
)

// TestPipelineQueueSimpleCase 通道
func TestPipelineQueueSimpleCase(t *testing.T) {
	done := make(chan any)
	defer close(done)

	zeros := take(done, repeat(done, 0, 1, 2, 3, 4), 3)
	short := sleep(done, 1*time.Second, zeros)
	long := sleep(done, 4*time.Second, short)
	pipeline := long
	// 会耗时13秒，因为等待的1秒钟会和4秒重叠，最终耗时 3*4秒+ 1秒
	for a := range pipeline {
		fmt.Printf("a value is %v \n", a)
	}
}

// TestPipelineBufQueueSimpleCase 增加buff, 队列的价值并不是减少了某个阶段的运行时间，而是减少了它处于阻塞状态的时间。
func TestPipelineBufQueueSimpleCase(t *testing.T) {
	done := make(chan any)
	defer close(done)

	zeros := take(done, repeat(done, 0, 1, 2, 3, 4), 3)
	bufferzeros := buffer(done, 2, zeros)

	short := sleep(done, 1*time.Second, bufferzeros)
	buffer := buffer(done, 2, short)
	long := sleep(done, 4*time.Second, buffer)
	pipeline := long
	// 还是会耗时13秒，因为等待的1秒钟会和4秒重叠，最终耗时 3*4秒+ 1秒
	for a := range pipeline {
		fmt.Printf("a value is %v \n", a)
	}
}
