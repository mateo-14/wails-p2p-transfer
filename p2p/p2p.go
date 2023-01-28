package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

type P2P struct {
	host host.Host
}

const (
	ProtocolPing = "/ping/1.0.0"
)

func New(privk crypto.PrivKey) (*P2P, error) {
	p := &P2P{}

	err := p.start(privk)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *P2P) start(privk crypto.PrivKey) error {
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/4000"), libp2p.Identity(privk))
	p.host = host
	host.SetStreamHandler(ProtocolPing, p.pingHandler)
	return err
}

func (p *P2P) ChangeKey(privk crypto.PrivKey) error {
	p.host.Close()
	return p.start(privk)
}

func (p *P2P) Addrs() []multiaddr.Multiaddr {
	return p.host.Addrs()
}

func (p *P2P) Id() string {
	return p.host.ID().String()
}

func (p *P2P) Connect(ctx context.Context, addr string) error {
	addrInfo, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(addrInfo)
	if err != nil {
		return err
	}

	err = p.host.Connect(ctx, *peerInfo)
	if err != nil {
		return err
	}

	p.host.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)

	s, err := p.host.NewStream(ctx, peerInfo.ID, ProtocolPing)
	if err != nil {
		return err
	}

	p.pingHandler(s)
	return nil
}

func (p *P2P) pingHandler(s network.Stream) {
	go func() {
		for {
			time.Sleep(time.Second)
			n, err := s.Write([]byte{1})
			if err != nil {
				fmt.Println("Error writing to stream: ", err)
			}

			fmt.Printf("Wrote %d bytes to stream\n", n)
		}
	}()
}

/* func (*P2P) streamHandler(s network.Stream) {
	go writeCounter(s)
	go readCounter(s)
}

func writeCounter(s network.Stream) {

	for {
		<-time.After(time.Second)


		n, err := s.Write([]byte("ping"))
		if err != nil {
			panic(err)
		}
	}
}

func readCounter(s network.Stream) {
	for {
		var counter uint64

		err := binary.Read(s, binary.BigEndian, &counter)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Received %d from %s\n", counter, s.ID())
	}
}
*/
