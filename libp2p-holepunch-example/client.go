package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	basichost "github.com/libp2p/go-libp2p/p2p/host/basic"
	manet "github.com/multiformats/go-multiaddr/net"
)

func runClient(ctx context.Context, role Role, serverAddrInfo peer.AddrInfo) {
	if !(role == Listener || role == Dialer) {
		panic(role)
	}

	// Create Host.
	hostOptions := getClientHostOptions(serverAddrInfo)
	host, err := libp2p.New(hostOptions...)
	if err != nil {
		log.Fatalf("Error creating host: %v\n", err)
	}
	defer host.Close()
	log.Printf("Host ID: %v\n", host.ID())

	for {
		time.Sleep(3 * time.Second)
		log.Printf("Connecting to server %v...\n", serverAddrInfo)
		err := host.Connect(ctx, serverAddrInfo)
		if err != nil {
			log.Printf("Error connecting to server: %v\n", err)
		} else {
			log.Println("Connected to server.")
			break
		}
	}

	idService := host.(*autorelay.AutoRelayHost).Host.(*basichost.BasicHost).IDService()
	for {
		hasPublic := false
		for _, addr := range idService.OwnObservedAddrs() {
			if manet.IsPublicAddr(addr) {
				hasPublic = true
				break
			}
		}
		if hasPublic {
			log.Printf("Observed self Addrs: %v\n", idService.OwnObservedAddrs())
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Connect to DHT.
	dhtNode, err := dht.New(ctx, host, dht.Mode(dht.ModeClient))
	if err != nil {
		log.Fatalf("Error creating DHT node: %v\n", err)
	}
	defer dhtNode.Close()
	err = dhtNode.Bootstrap(ctx)
	if err != nil {
		log.Fatalf("Error bootstraping DHT node: %v\n", err)
	}
	discovery := routing.NewRoutingDiscovery(dhtNode)

	if role == Listener {
		// Read messages from dialer.
		host.SetStreamHandler(protocol.TestingID, func(stream network.Stream) {
			defer stream.Close()
			log.Printf("Received stream from %v\n", stream.Conn().RemoteMultiaddr())
			reader := bufio.NewReader(stream)
			for {
				str, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						log.Printf("Error reading message from stream: %v\n", err)
					}
					return
				}
				str = strings.TrimSpace(str)
				log.Printf("Read message from stream: %s\n", str)
			}
		})

		// Advertise self to DHT periodically.
		for {
			time.Sleep(3 * time.Second)
			log.Println("Advertising to DHT...")
			_, err := discovery.Advertise(ctx, string(Listener))
			if err != nil {
				log.Printf("Error advertising to DHT: %v\n", err)
			} else {
				log.Println("Advertised to DHT.")
				time.Sleep(60 * time.Second)
			}
		}
	} else {
		for {
			// Discover listener
			var listenerAddrInfo peer.AddrInfo
			for {
				time.Sleep(3 * time.Second)
				log.Println("Finding listener from DHT...")
				peers, err := discovery.FindPeers(ctx, string(Listener))
				if err != nil {
					log.Printf("Error finding listener from DHT: %v\n", err)
				} else {
					listenerAddrInfo = <-peers
					if listenerAddrInfo.ID.Size() > 0 {
						break
					}
				}
			}
			log.Printf("Found listener: %v\n", listenerAddrInfo)

			// Send messages to listener.
			log.Println("Connecting to listener...")
			err = host.Connect(ctx, listenerAddrInfo)
			if err != nil {
				log.Printf("Error connecting to listener: %v\n", err)
			} else {
				log.Println("Connected to listener.")
				for {
					time.Sleep(1 * time.Second)
					log.Println("Creating stream...")
					stream, err := host.NewStream(ctx, listenerAddrInfo.ID, protocol.TestingID)
					if err != nil {
						log.Printf("Error creating stream: %v\n", err)
					} else {
						log.Printf("Created stream to %v\n", stream.Conn().RemoteMultiaddr())
						defer stream.Close()
						writer := bufio.NewWriter(stream)
						for {
							time.Sleep(1 * time.Second)
							fmt.Println("Sending message to listener...")
							_, err := writer.WriteString("Hello from dialer.\n")
							if err == nil {
								err = writer.Flush()
							}
							if err != nil {
								log.Printf("Error sending message to listener: %v\n", err)
								continue
							}
							fmt.Println("Sent message to listener.")
						}
					}
				}
			}
		}
	}
}
