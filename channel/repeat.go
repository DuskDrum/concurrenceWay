package channel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// repeat 一直重复，直到done消息传递进来，告诉他要停止。 使用done通道来传递关闭信息，这样可以避免go routine 内存泄露
func repeat(done <-chan any, values ...any) <-chan any {
	valueStream := make(chan interface{})
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

// take 从入参valueStream中去取值，
func take(done <-chan any, valueStream <-chan any, num int) <-chan any {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

// repeatFn 一直重复调用函数，直到done消息传递进来，告诉他要停止
func repeatFn(done <-chan any, fn func() interface{}) <-chan any {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

func toString(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}

		}
	}()
	return stringStream
}

func toInt(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	stringStream := make(chan int)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(int):
			}

		}
	}()
	return stringStream
}

func primeFinder(done <-chan interface{}, intStream <-chan int) <-chan any {
	primeStream := make(chan any)
	go func() {
		defer close(primeStream)
		for integer := range intStream {
			integer -= 1
			prime := true
			for divisor := integer - 1; divisor > 1; divisor-- {
				if integer%divisor == 0 {
					prime = false
					break
				}
			}

			if prime {
				select {
				case <-done:
					return
				case primeStream <- integer:
				}
			}
		}
	}()
	return primeStream
}

// fanIn模式
func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplexed := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplexed(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

// orDone or-done-channel 该通道在任何组件通道关闭时关闭
func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				// 这里的ok指的是 c通道是否已关闭，所以这里判断了c通道关闭，那么此goroutine也要关闭
				if ok == false {
					fmt.Println("通道已经关闭，向上传递此操作...")
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// tee-channel模式，把in拆分成两个
func tee(done <-chan any, in <-chan any) (<-chan any, <-chan any) {
	out1 := make(chan any)
	out2 := make(chan any)

	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			select {
			case out1 <- val:
				fmt.Printf("send to out1 %v \n", val)
			}
			select {
			case out2 <- val:
				fmt.Printf("send to out2 %v \n", val)
			}
		}
	}()
	return out1, out2
}

func bridge(done <-chan any, chanStream <-chan <-chan any) <-chan any {
	valStream := make(chan interface{}) // 1
	go func() {
		defer close(valStream)
		for {
			var stream <-chan any
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func genVals() <-chan <-chan any {
	chanStream := make(chan (<-chan any))
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			stream := make(chan any, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()

	return chanStream
}

func sleep(done <-chan any, d time.Duration, chanStream <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case val, ok := <-chanStream:
				if ok {
					fmt.Printf("收到消息, %v \n", val)
				} else {
					return
				}
				select {
				// 在这里睡会
				case <-time.After(d):
					valStream <- val
				}
			}
		}
	}()
	return valStream
}

func buffer(done <-chan any, bufSize int, chanStream <-chan any) <-chan any {
	bufStream := make(chan any, bufSize)
	go func() {
		defer close(bufStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-chanStream:
				if ok {
					bufStream <- v
				} else {
					return
				}
			}
		}
	}()
	return bufStream
}

func printGreeting(done <-chan any) error {
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}
func printFarewell(done <-chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}
func genGreeting(done <-chan any) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}
func genFarewell(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}
func locale(done <-chan any) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	// 等待五秒
	case <-time.After(5 * time.Second):
	}
	return "EN/US", nil
}

func printGreetingC(ctx context.Context) error {
	greeting, err := genGreetingC(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}
func printFarewellC(ctx context.Context) error {
	farewell, err := genFarewellC(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}
func genGreetingC(ctx context.Context) (string, error) {
	// 这里限制了一分钟，所以会报错 context deadline exceeded
	// 紧接着会影响到同一个上下文的 genFarewellC
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	switch locale, err := localeC(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}
func genFarewellC(ctx context.Context) (string, error) {
	switch locale, err := localeC(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func localeC(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
		// 这里会等待一分钟，但是ctx一秒就过期
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
