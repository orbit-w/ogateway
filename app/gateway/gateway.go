package gateway

import (
	gnetwork "github.com/orbit-w/golib/modules/net/network"
	"github.com/orbit-w/ogateway/app/oconfig"
	"github.com/spf13/viper"
)

/*
   @Author: orbit-w
   @File: gateway
   @2024 3月 周日 17:54
*/

func Serve() (IServer, error) {
	p := oconfig.Protocol()
	host := viper.GetString(oconfig.TagHost)

	protocol := parseProtocol(p)
	factory := getFactory(protocol)
	s := factory()
	if err := s.Serve(host); err != nil {
		return nil, err
	}

	return s, nil
}

type IServer interface {
	Serve(addr string) error
	Stop() error
}

type Factory func() IServer

var factories = make(map[gnetwork.Protocol]Factory)

func regFactory(name gnetwork.Protocol, f Factory) {
	factories[name] = f
}

func getFactory(name gnetwork.Protocol) Factory {
	return factories[name]
}

func parseProtocol(p string) gnetwork.Protocol {
	switch p {
	case "tcp":
		return gnetwork.TCP
	case "udp":
		return gnetwork.UDP
	case "kcp":
		return gnetwork.KCP
	default:
		return gnetwork.TCP
	}
}
