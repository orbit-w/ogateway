package agent

import (
	"context"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/ogateway/app/net/codec"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"net"
)

/*
   @Author: orbit-w
   @File: stream
   @2024 4月 周六 21:27
*/

type Stream struct {
	conn  net.Conn
	cli   client.XClient
	codec *codec.Codec
}

func NewStream() *Stream {
	return &Stream{
		codec: codec.NewCodec(MaxInPacketSize, false),
	}
}

func (s *Stream) Dial() error {
	addr := "localhost:8972"
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	xClient := client.NewXClient(share.StreamServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	stream, err := xClient.Stream(context.Background(), make(map[string]string))
	if err != nil {
		return err
	}
	s.cli = xClient
	s.conn = stream
	return nil
}

func (s *Stream) Recv(head, body []byte) (packet.IPacket, error) {
	return s.codec.BlockDecodeBody(s.conn, head, body)
}

func (s *Stream) Send(out packet.IPacket) error {
	body, err := s.codec.EncodeBody(out, false)
	if err != nil {
		return err
	}
	defer body.Return()
	_, err = s.conn.Write(body.Data())
	return err
}

func (s *Stream) Close() error {
	if s.conn != nil {
		_ = s.conn.Close()
	}
	if s.cli != nil {
		_ = s.cli.Close()
	}
	return nil
}
