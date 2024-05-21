package lib

import (
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/multiformats/go-multiaddr"
	glog "github.com/omgolab/go-commons/pkg/log"
)

type cfg struct {
	port   string
	logger glog.Logger
}

func getDefaultConfig() cfg {
	l, _ := glog.New()
	return cfg{
		port:   "9000",
		logger: l,
	}
}

type Option func(cfg *cfg) (libp2p.Option, error)

// add port
// default port is 9000
func WithPort(port int) Option {
	return func(cfg *cfg) (libp2p.Option, error) {
		if port < 0 || port > 65535 {
			return nil, errors.New("invalid port range")
		}
		return getListenAddresses(fmt.Sprintf("%d", port))
	}
}

// add logger
func WithLogger(log glog.Logger) Option {
	return func(cfg *cfg) (libp2p.Option, error) {
		if log == nil {
			return nil, errors.New("invalid logger")
		}
		cfg.logger = log
		return nil, nil
	}
}

// apply applies the given options to the config, returning the first error
// encountered (if any).
func (cfg *cfg) apply(opts ...Option) ([]libp2p.Option, error) {
	pOpts := make([]libp2p.Option, len(opts))
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		o, err := opt(cfg)
		if err != nil {
			return nil, err
		}

		if o == nil {
			continue
		}

		pOpts = append(pOpts, o)
	}
	return pOpts, nil
}

// getListenAddresses configures libp2p to use default listen address with the given port.
func getListenAddresses(port string) (libp2p.Option, error) {
	addrs := []string{
		"/ip4/0.0.0.0/tcp/" + port,
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1/webtransport",
		"/ip6/::/tcp/" + port,
		"/ip6/::/udp/" + port + "/quic-v1",
		"/ip6/::/udp/" + port + "/quic-v1/webtransport",
	}
	listenAddrs := make([]multiaddr.Multiaddr, 0, len(addrs))
	for _, s := range addrs {
		addr, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			return nil, err
		}
		listenAddrs = append(listenAddrs, addr)
	}

	return libp2p.ListenAddrs(listenAddrs...), nil
}
