package pool

import (
	"fmt"
	"sync"
	"testing"
)

// TestSimpleCase Pool模式的并发安全实现。
// Pool的主接口是它的Get方法。当调用时，Get将首先检查池中是否有可用的实例返回给调用者，如果没有，调用它的new方法来创建一个新实例。当完成时，调用者调用Put方法把工作的实例归还到池中，以供其他进程使用。
func TestSimpleCase(t *testing.T) {

	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance. ")
			return struct{}{}
		},
	}
	// 会调用New
	myPool.Get()
	// 这句也会调用New，因为Pool中没有实例，上一句的Get并没有Put回去
	instance := myPool.Get()
	// 正常的操作是Get后加一句 defer Put
	myPool.Put(instance)
	// Put回池中再Get，就无需再New了
	myPool.Get()
	// 这一句会继续New，因为上一句Get了出去没Put还回去
	myPool.Get()
}
