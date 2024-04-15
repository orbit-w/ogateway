package agent

import (
	"github.com/orbit-w/golib/bases/misc/utils"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/golib/modules/net/agent_stream"
	"github.com/orbit-w/ogateway/app/logger"
	"github.com/orbit-w/ogateway/app/net/onet"
	"go.uber.org/zap"
	"net"
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
	Authed     bool
	Idx        uint64 //increment id
	Uuid       int64  //agent id
	remoteAddr string
	state      atomic.Uint32
	conn       net.Conn
	cli        agent_stream.IStreamClient
	sender     ISender
	stream     agent_stream.IStream
}

func NewAgent(_Idx uint64, _conn net.Conn) *Agent {
	return &Agent{
		Idx:        _Idx,
		conn:       _conn,
		remoteAddr: _conn.RemoteAddr().String(),
	}
}

func (a *Agent) BindSender(s ISender) {
	a.sender = s
}

func (a *Agent) Proxy(out []byte) error {
	if !a.Authed {
		//登陆验证第一个消息包
		if err := a.auth(); err != nil {
			return err
		}
		if err := a.dial(); err != nil {
			return err
		}
		a.Authed = true
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
	cli := agent_stream.NewClient(agentStreamAddr)
	stream, err := cli.Stream()
	if err != nil {
		return err
	}
	a.stream = stream
	go a.handleLoop()
	return nil
}

func (a *Agent) auth() error {
	logger.ZLogger().Info("[Agent] authed", zap.Uint64("agent_id", a.Idx), zap.String("remote_addr", a.conn.RemoteAddr().String()))
	return nil
}

func (a *Agent) handleLoop() {
	var (
		err error
		in  []byte
	)

	utils.RecoverPanic()
	defer func() {
		a.safeReturn(err)
	}()

	for {
		in, err = a.stream.Recv()
		if err != nil {
			return
		}

		r := packet.Reader(in)
		if err = a.handleRespMsg(r); err != nil {
			return
		}
	}
}

func (a *Agent) handleRespMsg(in packet.IPacket) error {
	defer in.Return()
	p, err := in.ReadInt8()
	if err != nil {
		return err
	}

	switch p {
	case PatternNone:
		return a.sender.Send(in.Remain())
	case PatternKick:
		return a.Close()
	default:
		return AgentDecodePatternErr(p)
	}
}

func (a *Agent) safeReturn(err error) {
	if a.stream != nil {
		_ = a.stream.Close()
	}
	if a.conn != nil {
		_ = a.conn.Close()
	}
	if err != nil {
		switch {
		case onet.IsEOFError(err),
			onet.IsClosedConnError(err),
			onet.IsCancelError(err):
			//连接正常断开
			logger.ZLogger().Info("[Agent] connection closed", zap.Uint64("agent_id", a.Idx), zap.String("remote_addr", a.conn.RemoteAddr().String()))
		default:
			logger.ZLogger().Error("[Agent] abnormal connection disconnection", zap.Error(err))
		}
	}
}
