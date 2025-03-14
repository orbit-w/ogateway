package test

import (
	"testing"

	"github.com/orbit-w/ogateway/app"
	"github.com/orbit-w/ogateway/app/oconfig"
)

func Setup() {
	app.Run()
}

func Test_Main(t *testing.T) {
	oconfig.ParseConfig("../configs/config.toml")
	Setup()
}
