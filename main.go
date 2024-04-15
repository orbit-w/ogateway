package main

import (
	"flag"
	"github.com/orbit-w/ogateway/app/gateway"
	"github.com/orbit-w/ogateway/app/oconfig"
	"os"
	"os/signal"
	"syscall"
)

/*
   @Author: orbit-w
   @File: main
   @2024 3月 周三 23:50
*/

var configPath = flag.String("config", ".", "config file path")

func main() {
	flag.Parse()

	oconfig.ParseConfig(*configPath)

	stopper, err := gateway.Serve()
	if err != nil {
		panic(err)
	}
	defer func() {
		stopper()
	}()

	// Wait for exit signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
