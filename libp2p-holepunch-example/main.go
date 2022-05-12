package main

import (
	"context"
	"fmt"
	"log"

	logging "github.com/ipfs/go-log"
	"github.com/kelseyhightower/envconfig"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	logging.SetLogLevel("p2p-holepunch", "DEBUG")
	identify.ActivationThresh = 1

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("Error loading Config: %v\n", err)
	}
	log.Printf("Role: %v\n", config.Role)

	ctx := context.Background()

	serverID, serverIdentity := getIdentity(0)

	if config.Role == Server {
		runServer(ctx, serverIdentity)
	} else if config.Role == NAT {
		runNAT(ctx, config.NATRange, config.ExternalAddress)
	} else {
		if !(config.Role == Dialer || config.Role == Listener) {
			panic(config.Role)
		}
		configureGateway(config.GatewayAddress)

		// Get server AddrInfo.
		addrString := fmt.Sprintf("/ip4/%s/udp/%d/quic/", config.ServerAddress, defaultPort)
		addr, err := multiaddr.NewMultiaddr(addrString)
		if err != nil {
			panic(err)
		}
		serverAddrInfo := peer.AddrInfo{
			ID:    serverID,
			Addrs: []multiaddr.Multiaddr{addr},
		}
		log.Printf("Server AddrInfo: %v\n", serverAddrInfo)
		runClient(ctx, config.Role, serverAddrInfo, config.ListenerAddress)
	}

	log.Println("Exiting...")
}
