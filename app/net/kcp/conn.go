package okcp

import (
	"context"
	"github.com/orbit-w/golib/bases/misc/utils"
	"github.com/orbit-w/golib/bases/packet"
	gnetwork "github.com/orbit-w/golib/modules/net/network"
	"github.com/orbit-w/ogateway/app/logger"
	"github.com/orbit-w/ogateway/app/net/onet"
	"go.uber.org/zap"
	"io"
	"net"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: conn
   @2024 4月 周五 22:52
*/

type IAgent interface {
	Close() error
	Proxy(out []byte) error
}

type KcpConn struct {
	state  atomic.Uint32
	conn   net.Conn
	codec  *gnetwork.Codec
	ctx    context.Context
	cancel context.CancelFunc
	agent  IAgent
}

type ConnOptions struct {
	IsGzip            bool
	MaxIncomingPacket uint32
	Ctx               context.Context
	ReadTimeout       time.Duration
}

func NewKcpConn(_conn net.Conn, _agent IAgent, op ConnOptions) *KcpConn {
	ctx := op.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	cCtx, cancel := context.WithCancel(ctx)
	kc := &KcpConn{
		conn:   _conn,
		codec:  gnetwork.NewCodec(op.MaxIncomingPacket, op.IsGzip, op.ReadTimeout),
		ctx:    cCtx,
		cancel: cancel,
		agent:  _agent,
	}

	kc.state.Store(StatusNormal)
	return kc
}

func (kc *KcpConn) Send(data []byte) error {
	out, err := kc.codec.EncodeBodyRaw(data)
	if err != nil {
		return err
	}
	defer out.Return()
	if err = kc.conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
		return err
	}
	_, err = kc.conn.Write(out.Data())
	return err
}

func (kc *KcpConn) Close() error {
	if kc.state.CompareAndSwap(StatusNormal, StatusClosed) {
		if kc.agent != nil {
			_ = kc.agent.Close()
		}
	}
	return nil
}

func (kc *KcpConn) OnClose() error {
	if kc.state.CompareAndSwap(StatusNormal, StatusClosed) {
		if kc.conn != nil {
			_ = kc.conn.Close()
		}
	}
	return nil
}

func (kc *KcpConn) HandleLoop(head, body []byte) {
	var (
		err  error
		data packet.IPacket
	)

	defer utils.RecoverPanic()
	defer func() {
		_ = kc.Close()
		if err != nil {
			if err == io.EOF || onet.IsClosedConnError(err) {
				//连接正常断开
				logger.ZLogger().Info("[KcpConn] connection disconnected", zap.String("remote_addr", kc.conn.RemoteAddr().String()))
			} else {
				logger.ZLogger().Error("[KcpConn] abnormal connection disconnection", zap.Error(err), zap.String("remote_addr", kc.conn.RemoteAddr().String()))
			}
		}
	}()

	for {
		data, err = kc.codec.BlockDecodeBody(kc.conn, head, body)
		if err != nil {
			return
		}
		if err = kc.OnData(data); err != nil {
			return
		}
	}
}

func (kc *KcpConn) OnData(data packet.IPacket) error {
	defer data.Return()
	for len(data.Remain()) > 0 {
		if bytes, err := data.ReadBytes32(); err == nil {
			_ = kc.agent.Proxy(bytes)
		}
	}
	return nil
}
