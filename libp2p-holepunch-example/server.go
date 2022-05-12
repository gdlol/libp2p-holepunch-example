package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func runServer(ctx context.Context, identity libp2p.Option) {
	// Create host.
	hostOptions := getServerHostOptions(identity)
	host, err := libp2p.New(hostOptions...)
	if err != nil {
		log.Fatalf("Error creating host: %v\n", err)
	}
	defer host.Close()
	log.Printf("Host ID: %v\n", host.ID())

	// Create DHT node.
	dhtNode, err := dht.New(ctx, host, dht.Mode(dht.ModeServer))
	if err != nil {
		log.Fatalf("Error creating DHT node: %v\n", err)
	}
	defer dhtNode.Close()
	err = dhtNode.Bootstrap(ctx)
	if err != nil {
		log.Fatalf("Error bootstraping DHT node: %v\n", err)
	}

	log.Println("Server ready.")
	<-ctx.Done()
}
