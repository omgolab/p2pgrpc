package lib

import (
	"context"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/multiformats/go-multiaddr"
	glog "github.com/omgolab/go-commons/pkg/log"
)

// bootstrap is an optional helper to connect to the given peers and bootstrap
// the Peer DHT (and Bitswap). This is a best-effort function. Errors are only
// logged and a warning is printed when less than half of the given peers
// could be contacted. It is fine to pass a list where some peers will not be
// reachable.
func bootstrap(ctx context.Context, d *dht.IpfsDHT, h host.Host, log glog.Logger) {
	connected := make(chan struct{})
	peers, _ := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)

	var wg sync.WaitGroup
	for _, pInfo := range peers {
		wg.Add(1)
		go func(pInfo peer.AddrInfo) {
			defer wg.Done()
			err := h.Connect(ctx, pInfo)
			if err != nil {
				log.Warn(err.Error())
				return
			}
			log.Printf("Connected to", pInfo.ID)
			connected <- struct{}{}
		}(pInfo)
	}

	go func() {
		wg.Wait()
		close(connected)
	}()

	i := 0
	for range connected {
		i++
	}
	if nPeers := len(peers); i < nPeers/2 {
		log.Printf("only connected to %d bootstrap peers out of %d", i, nPeers)
	}

	err := d.Bootstrap(ctx)
	if err != nil {
		log.Error("dht bootstrap failed - ", err)
		return
	}
}

// createPeer creates a new libp2p Host with default settings.
func createPeer(ctx context.Context, opts []libp2p.Option, logger glog.Logger) host.Host {
	var dhtInstance *dht.IpfsDHT

	// Configure libp2p options
	options := []libp2p.Option{
		libp2p.Defaults,
		libp2p.NATPortMap(),       // Use NAT manager
		libp2p.EnableNATService(), // AutoNAT service
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error
			dhtInstance, err = dht.New(ctx, h, dht.Mode(dht.ModeAuto))
			return dhtInstance, err
		}),
		libp2p.EnableHolePunching(), // Enable Hole punching
	}
	options = append(options, opts...)

	// Create a new libp2p Host
	h, err := libp2p.New(options...)
	if err != nil {
		log.Fatalf("Failed to create libp2p host: %s", err)
	}

	// Bootstrap the DHT
	bootstrap(ctx, dhtInstance, h, logger)

	return h
}

// StartHost starts the libp2p host.
func StartHost(ctx context.Context, opts ...Option) error {
	log.Println("Launching host")
	cfg := getDefaultConfig()

	o, err := cfg.apply(opts...)
	if err != nil {
		return err
	}
	host := createPeer(ctx, o, cfg.logger)

	log.Printf("your hosts ID is: %s\n", host.ID().String())
	for _, addr := range host.Addrs() {
		ipfsAddr, err := multiaddr.NewMultiaddr("/ipfs/" + host.ID().String())
		if err != nil {
			panic(err)
		}
		peerAddr := addr.Encapsulate(ipfsAddr)
		log.Printf("host is listening on %s\n", peerAddr)
	}

	// rpcHost := gorpc.NewServer(host, lib.ProtocolID)

	// svc := PingService{}
	// err := rpcHost.Register(&svc)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Done")

	// for {
	// 	time.Sleep(time.Second * 1)
	// }

	return nil
}
