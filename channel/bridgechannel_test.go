package channel

import (
	"fmt"
	"testing"
)

// TestBridgeChannelSimpleCase 桥接channel模式
func TestBridgeChannelSimpleCase(t *testing.T) {
	//
	for v := range bridge(nil, genVals()) {
		fmt.Printf("从brideg中拿到的值 %v \n", v)
	}

}
