package gateway

import (
	"flag"
	"fmt"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/golib/modules/net/agent_stream"
	gnetwork "github.com/orbit-w/golib/modules/net/network"
	"github.com/orbit-w/ogateway/app/gateway/agent"
	"github.com/orbit-w/ogateway/app/net/onet"
	"github.com/orbit-w/ogateway/app/oconfig"
	"github.com/stretchr/testify/assert"
	"github.com/xtaci/kcp-go"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

/*
   @Author: orbit-w
   @File: gateway_test
   @2024 3月 周日 20:46
*/

var (
	once             sync.Once
	gs               *agent_stream.Server
	configPath       = flag.String("config", "../../configs", "config file path")
	streamServerHost = "127.0.0.1:8950"
)

func RunAgentStreamServer(handle func(stream agent_stream.IStream) error) {
	once.Do(func() {
		gs = new(agent_stream.Server)
		if err := gs.Serve(streamServerHost, handle); err != nil {
			panic(err)
		}
	})
}

func Test_Run(t *testing.T) {
	flag.Parse()
	oconfig.ParseConfig(*configPath)

	//启动 gateway server
	stopper, err := Serve()
	if err != nil {
		panic(err)
	}
	defer func() {
		stopper()
	}()

	//启动 agent_stream server
	RunAgentStreamServer(func(stream agent_stream.IStream) error {
		for {
			in, err := stream.Recv()
			if err != nil {
				break
			}
			fmt.Printf("agent_stream server recv: %s\n", string(in))
			w := packet.Writer()
			w.WriteInt8(agent.PatternNone)
			w.WriteString("hello, client")
			err = stream.Send(w.Data())
			if err != nil {
				t.Error(err)
			}
			w.Return()
		}
		return nil
	})

	//启动KCP客户端
	cli := NewKCPClient(t, "127.0.0.1")
	defer func() {
		_ = cli.Close()
	}()

	time.Sleep(time.Second * 30)
}

func Test_RunClient(t *testing.T) {
	NewKCPClient(t, "127.0.0.1")
	time.Sleep(time.Second * 5)
}

func NewKCPClient(t *testing.T, ip string) net.Conn {
	// 创建KCP客户端
	host := joinHostP(ip, "8900")
	conn, err := kcp.DialWithOptions(host, nil, 10, 3)
	assert.NoError(t, err)

	// 向服务器发送数据
	codec := gnetwork.NewCodec(gnetwork.MaxIncomingPacket, false, time.Second*60)

	// 从服务器读取数据
	var (
		head = make([]byte, 4)
		body = make([]byte, gnetwork.MaxIncomingPacket)
		in   packet.IPacket
	)
	go func() {
		for {
			in, err = codec.BlockDecodeBody(conn, head, body)
			if err != nil {
				if err == io.EOF || onet.IsClosedConnError(err) {
					fmt.Println("Server closed the connection")
					return
				}

				fmt.Printf("Error reading from server: %s", err.Error())
				break
			}

			ret, err := in.ReadBytes()
			assert.NoError(t, err)
			fmt.Println("Received from server: ", string(ret))
		}
	}()

	w := packet.Writer()
	w.WriteBytes32([]byte("Hello KCP Server!")) // 写入数据
	out, err := codec.EncodeBody(w)
	assert.NoError(t, err)
	_, err = conn.Write(out.Data())
	if err != nil {
		panic(err.Error())
	}
	assert.NoError(t, err)
	return conn
}
