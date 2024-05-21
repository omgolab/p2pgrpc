package lib

import (
	"context"
	"net"

	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	manet "github.com/multiformats/go-multiaddr/net"
)

type Libp2pListener struct {
	ctx        context.Context
	host       host.Host
	remotePeer peer.ID
	streams    chan net.Conn
}

func NewLibp2pListener(ctx context.Context, host host.Host, peerID peer.ID) *Libp2pListener {
	listener := &Libp2pListener{
		ctx:        ctx,
		host:       host,
		remotePeer: peerID,
		streams:    make(chan net.Conn),
	}
	host.SetStreamHandler(ProtocolID, listener.handleStream)
	return listener
}

func (l *Libp2pListener) handleStream(s network.Stream) {
	// Wrap the libp2p stream as a net.Conn object
	conn := wrapStream(s)
	l.streams <- conn
}

func (l *Libp2pListener) Accept() (net.Conn, error) {
	select {
	case <-l.ctx.Done():
		return nil, l.ctx.Err()
	case conn := <-l.streams:
		return conn, nil
	}
}

func (l *Libp2pListener) Close() error {
	l.host.RemoveStreamHandler(ProtocolID)
	close(l.streams)
	return nil
}

func (l *Libp2pListener) Addr() net.Addr {
	listenAddrs := l.host.Network().ListenAddresses()
	if len(listenAddrs) > 0 {
		for _, addr := range listenAddrs {
			if na, err := manet.ToNetAddr(addr); err == nil {
				return na
			}
		}
	}

	return defaultLocalAddr()
}
