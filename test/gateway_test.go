package test

import (
	"fmt"
	"github.com/orbit-w/ogateway/test/pb"
	"github.com/orbit-w/ogateway/test/pb/pb_core"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"testing"
	"time"

	"github.com/orbit-w/meteor/modules/net/packet"

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

func Test_RunKCPClient(t *testing.T) {

	conn := NewKCPClient(t, "47.120.6.89") //"47.120.6.89"
	time.Sleep(time.Second * 30)
	_ = conn.Close()
}

func Test_RunTCPClient(t *testing.T) {

	conn := NewTCPClient(t, "47.120.6.89") //"47.120.6.89"
	time.Sleep(time.Minute * 30)
	_ = conn.Close()
}

func NewTCPClient(t *testing.T, ip string) net.Conn {
	// 创建TCP客户端
	host := joinHostP(ip, "8900")
	conn, err := net.Dial("tcp", host)
	assert.NoError(t, err)

	run(t, conn)
	return conn
}

func NewKCPClient(t *testing.T, ip string) net.Conn {
	// 创建KCP客户端
	host := joinHostP(ip, "8900")
	conn, err := kcp.DialWithOptions(host, nil, 10, 3)
	assert.NoError(t, err)

	run(t, conn)
	return conn
}

func run(t *testing.T, conn net.Conn) {
	// 向服务器发送数据
	cliCodec := NewClientCodec()
	codec := gnetwork.NewCodec(gnetwork.MaxIncomingPacket, false, 0)

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
				rsp := new(pb_core.Request_SearchBook_Rsp)
				if err := proto.Unmarshal(msg.Data, rsp); err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Received message - PID: %d, Seq: %d, Data length: %d, Content: %s\n",
					msg.Pid, msg.Seq, len(msg.Data), rsp.Result.Content)
				// 这里可以根据不同的协议ID处理不同类型的消息
				// 例如: handleMessage(msg)
			}
		}
	}()

	req, err := proto.Marshal(&pb_core.Request_SearchBook{
		Query: "first",
	})
	assert.NoError(t, err)

	pack := cliCodec.Encode(req, 100, pb.PID_Core_Request_SearchBook)
	if err := conn.SetWriteDeadline(time.Now().Add(time.Second * 2)); err != nil {
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
}

func joinHostP(ip, port string) string {
	ipAddr := net.ParseIP(ip)
	return net.JoinHostPort(ipAddr.String(), port)
}
