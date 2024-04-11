package gateway

import (
	"flag"
	"fmt"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/golib/modules/net/agent_stream"
	gnetwork "github.com/orbit-w/golib/modules/net/network"
	"github.com/orbit-w/ogateway/app/net/onet"
	"github.com/orbit-w/ogateway/app/oconfig"
	"github.com/spf13/viper"
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
	once       sync.Once
	gs         *agent_stream.Server
	configPath = flag.String("config", "../../configs", "config file path")
)

func ServeTest(handle func(stream agent_stream.IStream) error) {
	once.Do(func() {
		gs = new(agent_stream.Server)
		if err := gs.Serve(viper.GetString(oconfig.TagHost), handle); err != nil {
			panic(err)
		}
	})
}

func Test_Run(t *testing.T) {
	flag.Parse()
	oconfig.ParseConfig(*configPath)

	server, err := Serve()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = server.Stop()
	}()

	ServeTest(func(stream agent_stream.IStream) error {
		for {
			in, err := stream.Recv()
			if err != nil {
				break
			}
			fmt.Printf("Server recv: %s", string(in))
			err = stream.Send([]byte("hello, client"))
			if err != nil {
				t.Error(err)
			}
		}
		return nil
	})

	cli := NewClient(t)
	defer func() {
		_ = cli.Close()
	}()

	time.Sleep(time.Second * 30)
}

func NewClient(t *testing.T) net.Conn {
	// 创建KCP客户端
	conn, err := kcp.DialWithOptions(viper.GetString(oconfig.TagHost), nil, 10, 3)
	assert.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

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
			}

			fmt.Printf("Received from server: %s\n", string(in.Data()))
		}
	}()

	out, err := codec.EncodeBodyRaw([]byte("Hello KCP Server!"))
	assert.NoError(t, err)
	_, err = conn.Write(out.Data())
	assert.NoError(t, err)
	return conn
}
