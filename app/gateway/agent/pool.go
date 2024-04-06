package agent

import "github.com/orbit-w/golib/core/network"

/*
   @Author: orbit-w
   @File: pool
   @2024 4月 周六 22:08
*/

var (
	headPool = network.NewBufferPool(headLen)
	bodyPool = network.NewBufferPool(MaxInPacketSize)
)
