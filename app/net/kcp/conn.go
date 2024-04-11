package okcp

import (
	"context"
	"fmt"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/ogateway/app/net/codec"
	"github.com/orbit-w/ogateway/app/net/onet"
	"io"
	"log"
	"net"
	"runtime/debug"
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
	codec  *codec.Codec
	ctx    context.Context
	cancel context.CancelFunc
	agent  IAgent
}

type ConnOptions struct {
	MaxIncomingPacket uint32
	Ctx               context.Context
}

func NewKcpConn(_conn net.Conn, _agent IAgent, op ConnOptions) *KcpConn {
	ctx := op.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	cCtx, cancel := context.WithCancel(ctx)
	kc := &KcpConn{
		conn:   _conn,
		codec:  codec.NewCodec(op.MaxIncomingPacket, false),
		ctx:    cCtx,
		cancel: cancel,
		agent:  _agent,
	}

	kc.state.Store(StatusNormal)
	return kc
}

func (kc *KcpConn) Send(data []byte) (err error) {
	if err = kc.conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
		return err
	}
	_, err = kc.conn.Write(data)
	return
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

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println("stack: ", string(debug.Stack()))
		}
		_ = kc.Close()
		if err != nil {
			if err == io.EOF || onet.IsClosedConnError(err) {
				//连接正常断开
			} else {
				log.Println(fmt.Errorf("[TcpServer] tcp_conn disconnected: %s", err.Error()))
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
