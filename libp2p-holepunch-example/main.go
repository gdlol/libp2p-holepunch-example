package main

import (
	"context"
	"fmt"
	"log"
	"os"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	logging.SetLogLevel("p2p-holepunch", "DEBUG")
	identify.ActivationThresh = 1

	// Get role of the node.
	roleEnv := os.Getenv("ROLE")
	var role Role = Role(roleEnv)
	switch role {
	case Server:
	case NAT:
	case Listener:
	case Dialer:
	default:
		log.Fatalf("Unknown role: %s\n", roleEnv)
	}
	log.Printf("Role: %v\n", role)

	ctx := context.Background()

	serverID, serverIdentity := getIdentity(0)

	// Run server.
	if role == Server {
		runServer(ctx, serverIdentity)
	} else {
		// Run NAT.
		if role == NAT {
			natRange := os.Getenv("NAT_RANGE")
			log.Printf("NAT range: %s\n", natRange)
			serverNetNATAddress := os.Getenv("SERVER_NET_NAT_ADDRESS")
			log.Printf("NAT address: %s\n", serverNetNATAddress)
			clientAddress := os.Getenv("CLIENT_ADDRESS")
			runNAT(ctx, natRange, serverNetNATAddress, clientAddress)
		} else {
			// Run Client.
			if !(role == Dialer || role == Listener) {
				panic(role)
			}
			serverAddress := os.Getenv("SERVER_ADDRESS")
			log.Printf("Server address: %s\n", serverAddress)

			// Configure routing.
			natAddress := os.Getenv("NAT_ADDRESS")
			log.Printf("NAT address: %s\n", natAddress)
			configureGateway(natAddress)

			// Get server AddrInfo.
			addrString := fmt.Sprintf("/ip4/%s/udp/%d/quic/", serverAddress, defaultPort)
			addr, err := multiaddr.NewMultiaddr(addrString)
			if err != nil {
				panic(err)
			}
			serverAddrInfo := peer.AddrInfo{
				ID:    serverID,
				Addrs: []multiaddr.Multiaddr{addr},
			}
			log.Printf("Server AddrInfo: %v\n", serverAddrInfo)
			runClient(ctx, role, serverAddrInfo)
		}
	}

	log.Println("Exiting...")
}
