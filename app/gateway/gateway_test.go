package gateway

import (
	"fmt"
	"testing"
)

/*
   @Author: orbit-w
   @File: gateway_test
   @2024 3月 周日 20:46
*/

func Test_Run(t *testing.T) {
	// Create a KCP listener
	fmt.Println(200 * 1024)
	fmt.Println((1 << 18) / 1024)
	fmt.Println(1<<18 - 1)
}
