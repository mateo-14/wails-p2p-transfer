package p2p

import (
	"context"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type P2P struct {
	host     host.Host
	Notifiee network.Notifiee
}

type HostData struct {
	Addr string `json:"address"`
	Id   string `json:"id"`
}

const ProtocolPing = "/ping/1.0.0"

func New(notifiee network.Notifiee) (*P2P, error) {
	privk, err := loadPrivateKey()
	if err != nil {
		privk, err = generatePrivateKey()
		if err != nil {
			return nil, err
		}

		err = storePrivateKey(privk)
		if err != nil {
			return nil, err
		}
	}

	p := &P2P{}
	if notifiee != nil {
		p.Notifiee = notifiee
	}

	err = p.start(privk)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *P2P) start(privk crypto.PrivKey) error {
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/4000"), libp2p.Identity(privk))

	if p.Notifiee != nil {
		fmt.Printf("%p\n", p.Notifiee)
		host.Network().Notify(p.Notifiee)
	}

	p.host = host
	// host.SetStreamHandler(ProtocolPing, p.pingHandler)
	for _, peer := range host.Peerstore().Peers() {
		fmt.Println(peer)
	}
	return err
}

func (p *P2P) RegenerateKey() error {
	p.host.Close()

	privk, err := generatePrivateKey()
	if err != nil {
		return err
	}

	err = storePrivateKey(privk)
	if err != nil {
		return err
	}

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

	/* _, err = p.host.NewStream(ctx, peerInfo.ID, ProtocolPing)
	if err != nil {
		return err
	} */

	// p.pingHandler(s)
	return nil
}

func (p *P2P) GetHostData() *HostData {
	return &HostData{
		Addr: p.Addrs()[0].String(),
		Id:   p.Id(),
	}
}

/* func (p *P2P) pingHandler(s network.Stream) {
	s.Conn().ID()
	go func() {
		for {
			time.Sleep(time.Second)
			n, err := s.Write([]byte{1})
			if err != nil {
				fmt.Println("Ping: Error writing to stream: ", err)
				return
			}

			fmt.Printf("Ping: Wrote %d bytes to stream\n", n)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			buf := make([]byte, 1)
			n, err := s.Read(buf)
			if err != nil {
				fmt.Println("Ping: Error reading from stream: ", err)
				return
			}

			fmt.Printf("Ping: Read %d bytes from stream \n", n)
		}
	}()
} */

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

func loadPrivateKey() (crypto.PrivKey, error) {
	file, err := os.ReadFile("private.key")
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(file)
}

func storePrivateKey(privk crypto.PrivKey) error {
	privkb, err := crypto.MarshalPrivateKey(privk)
	if err != nil {
		return err
	}

	return os.WriteFile("private.key", privkb, 0644)
}

func generatePrivateKey() (crypto.PrivKey, error) {
	privk, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	return privk, err
}
