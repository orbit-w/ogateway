package test

import (
	"fmt"
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

func Test_RunClient(t *testing.T) {
	conn := NewKCPClient(t, "127.0.0.1") //"47.120.6.89"
	time.Sleep(time.Second * 30)
	_ = conn.Close()
}

func NewKCPClient(t *testing.T, ip string) net.Conn {
	// 创建KCP客户端
	host := joinHostP(ip, "8900")
	conn, err := kcp.DialWithOptions(host, nil, 10, 3)
	assert.NoError(t, err)

	// 向服务器发送数据
	codec := gnetwork.NewCodec(gnetwork.MaxIncomingPacket, false, time.Second*60)

	// 从服务器读取数据
	//var (
	//	head = make([]byte, 4)
	//	body = make([]byte, gnetwork.MaxIncomingPacket)
	//	in   packet.IPacket
	//)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF || onet.IsClosedConnError(err) {
					fmt.Println("Server closed the connection")
					return
				}

				fmt.Printf("Error reading from server: %s\n", err.Error())
				fmt.Println(err.Error())
				fmt.Println("========")
				fmt.Println(n)
				fmt.Println("========")
				break
			}
		}
	}()

	w := packet.WriterP(1024)
	w.WriteBytes32([]byte("Hello KCP Server!")) // 写入数据
	out := codec.EncodeBody(w.Data(), 0)
	if err = conn.SetWriteDeadline(time.Now().Add(time.Second * 2)); err != nil {
		panic(err.Error())
	}
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
