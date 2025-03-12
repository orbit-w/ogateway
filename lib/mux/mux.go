package multiplexers

/*
   @Author: orbit-w
   @File: mux.go
   @2025 3月 周一 23:50
*/

import (
	"github.com/orbit-w/mux-go/multiplexers"
)

var mux *multiplexers.Multiplexers

func InitMultiplexers(host string) {
	mux = multiplexers.NewWithDefaultConf(host)
}

func Dial() (multiplexers.IConn, error) {
	return mux.Dial()
}
