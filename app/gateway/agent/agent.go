package agent

import (
	"fmt"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/golib/core/network"
	"github.com/orbit-w/ogateway/app/net/onet"
	"io"
	"log"
	"net"
	"runtime/debug"
	"sync/atomic"
)

/*
   @Author: orbit-w
   @File: agent
   @2024 4月 周二 22:30
*/

type ISender interface {
	Send(data []byte) (err error)
}

type IStream interface {
	Close() error
}

type Agent struct {
	Authed bool
	Idx    uint64 //increment id
	Uuid   int64  //agent id
	state  atomic.Uint32
	conn   net.Conn
	sender ISender
	stream *Stream
}

func NewAgent(_Idx uint64, _conn net.Conn) *Agent {
	return &Agent{
		Idx:  _Idx,
		conn: _conn,
	}
}

func (a *Agent) BindSender(s ISender) {
	a.sender = s
}

func (a *Agent) Send(out packet.IPacket) error {
	if !a.Authed {
		//登陆验证第一个消息包
		if err := a.auth(); err != nil {
			return err
		}
		if err := a.dial(); err != nil {
			return err
		}
		a.Authed = true
		return nil
	}
	return a.stream.Send(out)
}

func (a *Agent) Close() error {
	if a.state.CompareAndSwap(StatusNormal, StatusClosed) {
		if a.stream != nil {
			_ = a.stream.Close()
		}
		if a.conn != nil {
			_ = a.conn.Close()
		}
	}
	return nil
}

func (a *Agent) dial() error {
	a.stream = NewStream()
	if err := a.stream.Dial(); err != nil {
		return err
	}
	go a.handleLoop()
	return nil
}

func (a *Agent) auth() error {
	return nil
}

func (a *Agent) handleLoop() {
	var (
		err     error
		headBuf = headPool.Get().(*network.Buffer)
		bodyBuf = bodyPool.Get().(*network.Buffer)
		head    = headBuf.Bytes
		body    = bodyBuf.Bytes
		in      packet.IPacket
	)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println("stack: ", string(debug.Stack()))
		}

		_ = a.Close()
		headPool.Put(headBuf)
		bodyPool.Put(bodyBuf)
		if err != nil {
			if err == io.EOF || onet.IsClosedConnError(err) {
				//连接正常断开
			} else {
				log.Println(fmt.Errorf("[TcpServer] tcp_conn disconnected: %s", err.Error()))
			}
		}
	}()

	for {
		in, err = a.stream.Recv(head, body)
		if err != nil {
			return
		}

		if err = a.handleMsg(in); err != nil {
			return
		}
	}
}

func (a *Agent) handleMsg(in packet.IPacket) error {
	defer in.Return()
	p, err := in.ReadInt8()
	if err != nil {
		return err
	}

	switch p {
	case PatternKick:
		_ = a.Close()
		return nil
	default:
		return a.sender.Send(in.Remain())
	}
}
