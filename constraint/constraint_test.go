package constraint

import (
	"fmt"
	"testing"
)

// TestConstraintCase 超时机制，不需要
func TestConstraintCase(t *testing.T) {
	data := make([]int, 4)
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}

}
