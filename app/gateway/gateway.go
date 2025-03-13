package gateway

import (
	"context"
	"net"

	gnetwork "github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/ogateway/app/logger"
	"github.com/orbit-w/ogateway/app/oconfig"
	multiplexers "github.com/orbit-w/ogateway/lib/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

/*
   @Author: orbit-w
   @File: gateway
   @2024 3月 周日 17:54
*/

func Serve() (stopper func(ctx context.Context) error, err error) {

	host := joinHost()
	p := oconfig.Protocol()
	protocol := parseProtocol(p)
	factory := getFactory(protocol)
	s := factory()

	if err = s.Serve(host); err != nil {
		return nil, err
	}

	logger.ZLogger().Info("gateway listened...", zap.String("Port", viper.GetString(oconfig.TagPort)), zap.String("Protocol", p))

	stopper = func(ctx context.Context) error {
		multiplexers.CloseAll()

		if err = s.Stop(); err != nil {
			logger.ZLogger().Error("gateway stop error", zap.Error(err))
		}
		logger.StopLogger()
		return nil
	}

	return stopper, nil
}

func joinHost() string {
	ip := viper.GetString(oconfig.TagIp)
	port := viper.GetString(oconfig.TagPort)
	ipAddr := net.ParseIP(ip)
	return net.JoinHostPort(ipAddr.String(), port)
}

func joinHostP(ip, port string) string {
	ipAddr := net.ParseIP(ip)
	return net.JoinHostPort(ipAddr.String(), port)
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
