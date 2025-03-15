package test

import (
	"fmt"
	"github.com/orbit-w/meteor/modules/net/packet"
	"io"
	"net"
	"testing"
	"time"

	gnetwork "github.com/orbit-w/meteor/modules/net/network"
	"github.com/orbit-w/ogateway/app/net/onet"
	"github.com/stretchr/testify/assert"
	"github.com/xtaci/kcp-go"
)

/*
   @Author: orbit-w
   @File: gateway_test
   @2024 3月 周日 20:46
*/

func Test_RunClient(t *testing.T) {
	conn := NewKCPClient(t, "127.0.0.1") //"47.120.6.89"
	time.Sleep(time.Minute * 30)
	_ = conn.Close()
}

func NewKCPClient(t *testing.T, ip string) net.Conn {
	// 创建KCP客户端
	host := joinHostP(ip, "8900")
	conn, err := kcp.DialWithOptions(host, nil, 10, 3)
	assert.NoError(t, err)

	// 向服务器发送数据
	cliCodec := NewClientCodec()
	codec := gnetwork.NewCodec(gnetwork.MaxIncomingPacket, false, 0)

	// 从服务器读取数据
	//var (
	//	head = make([]byte, 4)
	//	body = make([]byte, gnetwork.MaxIncomingPacket)
	//	in   packet.IPacket
	//)
	go func() {
		head := make([]byte, 4)
		body := make([]byte, 4096)
		for {
			data, _, err := codec.BlockDecodeBody(conn, head, body)
			if err != nil {
				if err == io.EOF || onet.IsClosedConnError(err) {
					fmt.Println("Server closed the connection")
					return
				}

				fmt.Printf("Error reading from server: %s\n", err.Error())
				fmt.Println(err.Error())
				fmt.Println("========")

				break
			}

			// 使用ClientCodec解码接收到的消息
			messages, err := cliCodec.Decode(data, func(pid uint32) bool {
				// 这里可以根据协议ID判断是否需要读取seq
				// 返回true表示需要读取seq，返回false表示不需要
				return true
			})

			if err != nil {
				fmt.Printf("Error decoding message: %s\n", err.Error())
				continue
			}

			// 处理解码后的消息
			for _, msg := range messages {
				fmt.Printf("Received message - PID: %d, Seq: %d, Data length: %d, Content: %s\n",
					msg.Pid, msg.Seq, len(msg.Data), string(msg.Data))
				// 这里可以根据不同的协议ID处理不同类型的消息
				// 例如: handleMessage(msg)
			}
		}
	}()

	pack := cliCodec.Encode([]byte("Hello KCP Server!"), 100, 10)
	if err = conn.SetWriteDeadline(time.Now().Add(time.Second * 2)); err != nil {
		panic(err.Error())
	}

	w := packet.Writer(1024)
	w.WriteBytes32(pack.Data())

	out := codec.EncodeBody(w.Data(), 0)

	n, err := conn.Write(out.Data())
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(n)
	assert.NoError(t, err)
	return conn
}

func joinHostP(ip, port string) string {
	ipAddr := net.ParseIP(ip)
	return net.JoinHostPort(ipAddr.String(), port)
}
