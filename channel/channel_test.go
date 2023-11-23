package channel

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

// TestSimpleCase 一个 channel充当着信息传送的管道，值可以沿着channel传递，然后在下游读出、
// 当你使用channel时，你会将一个值传递给一个chan变量，然后你程序中的某个地方将它从channel中读出
// <-chan 代表只读通道
// chan<- 代表只发通道
func TestSimpleCase(t *testing.T) {
	stringStream := make(chan string)
	go func() {
		time.Sleep(10 * time.Second)
		stringStream <- "Hello channels!"
	}()
	// 这里会被阻塞，直到管道有值
	fmt.Println(<-stringStream)
}

// TestCloseChannelCase 从已经关闭了的channel中读取元素，会读取到0值
// channel关闭后，我们仍然可以继续在这个channel上执行读取操作，但是读取的内容是错误的，需要判断ok字段是否是false
func TestCloseChannelCase(t *testing.T) {
	stringStream := make(chan string)
	close(stringStream)
	str, ok := <-stringStream
	// string的零值是""
	fmt.Println(str)
	// ok也会是false
	fmt.Println(ok)
}

// TestRangeCase channel关闭时会自动中断循环
// 注意该循环不需要退出条件，并且range方能不返回第二个布尔值。处理一个已关闭的channel的细节可以让你保持循环简洁。
func TestRangeCase(t *testing.T) {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()
	for integer := range intStream {
		fmt.Printf("%v ", integer)
	}
}

// TestCloseNotifyCase 可以使用channel来作为通知。如果只是用channel来作为信号的(能从channel读到内容就行)，那么可以使用close来模拟notifyAll的功能
func TestCloseNotifyCase(t *testing.T) {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// 这里不关心从chan中拿到的是什么，拿到的是零值也可以，所以这里可以看做是一块挡板
			<-begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}
	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}

// TestBuffChannelCase 缓冲chan
func TestBuffChannelCase(t *testing.T) {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)
	// 容量为4的缓冲chan
	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, " sending:%d \n", i)
			intStream <- i
		}
	}()
	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "received %v. \n", integer)
	}
}

// TestSelectChannelSimpleCase 使用select来将多个chan绑定到一起，并且同时处理取消、超时、等待和默认值
func TestSelectChannelSimpleCase(t *testing.T) {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()
	fmt.Println("Blocking on read...")
	select {
	// 这一句会阻塞5秒，所以select会一直等待，知道从c取到值
	case <-c:
		fmt.Printf("Unblocked %v later. \n", time.Since(start))
	}
}

// TestSelectChannelDuplicateCase select中存在多个chan会怎么样
func TestSelectChannelDuplicateCase(t *testing.T) {
	cl := make(chan interface{})
	close(cl)
	c2 := make(chan interface{})
	close(c2)
	var c1Count, c2Count int
	for i := 1000; i > 0; i-- {
		// 可以看出来这里的c1和c2基本上是对半的概率，所以可以认为是随机选择
		select {
		case <-cl:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\n c2Count: %d \n",
		c1Count,
		c2Count)
}

// TestSelectChannelSimpleCase 使用select来将多个chan绑定到一起，并且同时处理取消、超时、等待和默认值
func TestSelectChannelDefaultCase(t *testing.T) {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()
	fmt.Println("Blocking on read...")
	select {
	// 这一句会阻塞5秒，所以select会一直等待，知道从c取到值
	case <-c:
		fmt.Printf("Unblocked %v later. \n", time.Since(start))
	// 由于有default，所以select并不会去等待从c中获取值，而是直接从这里拿到内容
	default:
		fmt.Printf("Unblocked %v later. \n", time.Since(start))
	}
}

// TestSelectChannelTimeoutCase 超时机制，不需要
func TestSelectChannelTimeoutCase(t *testing.T) {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()
	fmt.Println("Blocking on read...")
	select {
	// 这一句会阻塞5秒，所以select会一直等待，知道从c取到值
	case <-c:
		fmt.Printf("Unblocked %v later. \n", time.Since(start))
	// 1秒钟会超时
	case <-time.After(1 * time.Second):
		fmt.Printf("Timed out")
	}
}
