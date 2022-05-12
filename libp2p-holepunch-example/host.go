package main

import (
	"fmt"
	"math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"

	quic "github.com/libp2p/go-libp2p-quic-transport"
)

type randomReader struct {
	random *rand.Rand
}

func (reader randomReader) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		p[i] = byte(reader.random.Intn(256))
	}
	return len(p), nil
}

// Deterministic key.
func getIdentity(seed int64) (peer.ID, libp2p.Option) {
	reader := randomReader{
		random: rand.New(rand.NewSource(seed)),
	}
	privKey, _, err := crypto.GenerateEd25519Key(&reader)
	if err != nil {
		panic(err)
	}
	peerID, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		panic(err)
	}
	identity := libp2p.Identity(privKey)
	return peerID, identity
}

func getClientHostOptions(identity libp2p.Option, serverAddrInfo peer.AddrInfo) []libp2p.Option {
	return []libp2p.Option{
		identity,
		libp2p.Transport(quic.NewTransport),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", defaultPort)),
		libp2p.ForceReachabilityPrivate(),
		libp2p.EnableAutoRelay(),
		libp2p.StaticRelays([]peer.AddrInfo{serverAddrInfo}),
		libp2p.EnableHolePunching(),
	}
}

func getServerHostOptions(identity libp2p.Option) []libp2p.Option {
	return []libp2p.Option{
		identity,
		libp2p.Transport(quic.NewTransport),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", defaultPort)),
		libp2p.ForceReachabilityPublic(),
		libp2p.EnableRelayService(),
		libp2p.EnableNATService(),
	}
}
