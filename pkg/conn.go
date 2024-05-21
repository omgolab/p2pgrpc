package lib

import (
	"net"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	manet "github.com/multiformats/go-multiaddr/net"
)

type Libp2pConn struct {
	stream network.Stream
}

func wrapStream(stream network.Stream) net.Conn {
	return &Libp2pConn{stream: stream}
}

func (c *Libp2pConn) Read(b []byte) (int, error) {
	return c.stream.Read(b)
}

func (c *Libp2pConn) Write(b []byte) (int, error) {
	return c.stream.Write(b)
}

func (c *Libp2pConn) Close() error {
	return c.stream.Close()
}

func (c *Libp2pConn) LocalAddr() net.Addr {
	if addr := c.stream.Conn().LocalMultiaddr(); addr != nil {
		if ipAddr, err := manet.ToNetAddr(addr); err == nil {
			return ipAddr
		}
	}
	return defaultLocalAddr() // Return default if conversion fails
}

func (c *Libp2pConn) RemoteAddr() net.Addr {
	if addr := c.stream.Conn().RemoteMultiaddr(); addr != nil {
		if ipAddr, err := manet.ToNetAddr(addr); err == nil {
			return ipAddr
		}
	}
	return defaultLocalAddr() // Return default if conversion fails
}

func (c *Libp2pConn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *Libp2pConn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

func (c *Libp2pConn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}
