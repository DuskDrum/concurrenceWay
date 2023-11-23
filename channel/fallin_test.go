package channel

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestPrimeFinderCase(t *testing.T) {
	r := func() interface{} { return rand.Intn(500000000) }
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	// 生成随机数
	randIntStream := toInt(done, repeatFn(done, r))
	fmt.Println("Primes:")
	// 筛选里面的素数
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	// 速度回很慢，大概二十多秒，因为算法有问题
	fmt.Printf("Search took: %v", time.Since(start))
}

func TestPrimeFinderFallInCase(t *testing.T) {
	r := func() interface{} { return rand.Intn(500000000) }
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	// 生成随机数
	randIntStream := toInt(done, repeatFn(done, r))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)

	fmt.Println("Primes:")
	// 筛选里面的素数
	for i := 0; i < numFinders; i++ {
		// 每次都会make一个channel出来
		finders[i] = primeFinder(done, randIntStream)
	}

	// finders是一个channel的数组，fallIn把这些数组里的内容归拢到一个channel
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	// 速度慢，比原始要快很多，大概三秒
	fmt.Printf("Search took: %v", time.Since(start))
}
