package gateway

import (
	"context"
	"github.com/orbit-w/golib/modules/net/network"
	"github.com/orbit-w/ogateway/app/gateway/agent"
	"github.com/orbit-w/ogateway/app/logger"
	okcp "github.com/orbit-w/ogateway/app/net/kcp"
	"github.com/xtaci/kcp-go"
	"go.uber.org/zap"
	"net"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: server
   @2024 4月 周六 12:31
*/

func init() {
	regFactory(network.KCP, func() IServer {
		return new(KcpServer)
	})
}

type Stopper interface {
	Stop() error
}

// KcpServer 当你使用kcp的时候，你必须设置Timeout，利用timeout保持连接的检测。
// 因为kcp-go本身不提供keepalive/heartbeat的功能，当服务器宕机重启的时候，
// 原有的连接没有任何异常，只会hang住，我们只能依靠Timeout避免hang住
type KcpServer struct {
	idIncr  atomic.Uint64
	stopper Stopper
}

func (kcpS *KcpServer) Serve(addr string) error {
	// Create a KCP listener
	listener, err := kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		panic(err)
	}

	server := new(network.Server)
	server.Serve("kcp", listener, func(ctx context.Context, _conn net.Conn, maxIncomingPacket uint32, head, body []byte) {
		idx := kcpS.Idx()
		oAgent := agent.NewAgent(idx, _conn)

		conn := okcp.NewKcpConn(_conn, oAgent, okcp.ConnOptions{
			Ctx:               ctx,
			MaxIncomingPacket: maxIncomingPacket,
			ReadTimeout:       network.ReadTimeout,
		})
		defer func() {
			_ = conn.Close()
		}()

		oAgent.BindSender(conn)
		logger.ZLogger().Info("new kcp connection, binding agent", zap.Uint64("AgentId", idx), zap.String("RemoteAddr", _conn.RemoteAddr().String()))
		conn.HandleLoop(head, body)
	}, network.AcceptorOptions{
		MaxIncomingPacket: MaxInPacketSize,
		IsGzip:            false,
	})

	kcpS.stopper = server
	return nil
}

func (kcpS *KcpServer) Stop() error {
	return kcpS.stopper.Stop()
}

func (kcpS *KcpServer) Idx() uint64 {
	return kcpS.idIncr.Add(1)
}
