package waitgroup

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// TestSimpleClosureCase 闭包简单case，可以发现闭包可以从创建他们作用域中获取变量，在goroutine中运行一个闭包，闭包会直接使用原值的地址空间内执行
func TestSimpleClosureCase(t *testing.T) {
	var wg sync.WaitGroup
	salutation := "hello"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()
	wg.Wait()
	fmt.Println(salutation)
	assert.True(t, salutation == "welcome")
}

// TestCycleClosureCase 闭包循环case，字符串类型的切片进行循环闭包，结果每个循环得到的都是最后一位，这样设计的原因是因为goroutine中的代码不确定什么时候才能执行，所以可能会导致切片被改变了，从而导致其他类似于切片超出范围，内存被回收等问题
// 通常在任何goroutine执行前，循环就执行完了。goland也会有提示：Loop variables captured by 'func' literals in 'go' statements might have unexpected values
// 解决方案有两个
func TestCycleClosureCase(t *testing.T) {
	var wg sync.WaitGroup
	slice := []string{"hello", "greetings", "good day"}
	for _, salutation := range slice {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
			assert.True(t, salutation == "good day")
		}()
	}
	wg.Wait()
}

// TestCycleClosureSolver1Case 解决闭包循环问题，方案1：将salutation的副本传递到闭包中(例子中把i也传进去了，是为了校验assert，可以忽略)
func TestCycleClosureSolver1Case(t *testing.T) {
	var wg sync.WaitGroup
	slice := []string{"hello", "greetings", "good day"}
	for i, salutation := range slice {
		wg.Add(1)
		go func(salutation string, i int) {
			defer wg.Done()
			fmt.Println(salutation)
			assert.True(t, salutation == slice[i])
		}(salutation, i)
	}
	wg.Wait()
}

// TestCycleClosureSolver2Case 解决闭包循环问题，方案2：将join-point放到循环里，观察打印出来的内容
func TestCycleClosureSolver2Case(t *testing.T) {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
		wg.Wait()
	}
}
