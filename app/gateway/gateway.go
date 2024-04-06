package gateway

/*
   @Author: orbit-w
   @File: gateway
   @2024 3月 周日 17:54
*/

func Serve() {
}

type IServer interface {
	Serve(addr string) error
	Stop() error
}
