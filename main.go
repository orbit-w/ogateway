package main

import (
	"flag"

	"github.com/orbit-w/ogateway/app"
	"github.com/orbit-w/ogateway/app/oconfig"
)

/*
   @Author: orbit-w
   @File: main
   @2024 3月 周三 23:50
*/

var configPath = flag.String("config", "./configs/config.toml", "config file path")

func main() {
	flag.Parse()

	oconfig.ParseConfig(*configPath)

	app.Run()
}
