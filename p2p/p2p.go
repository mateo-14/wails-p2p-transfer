package p2p

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
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

type Message struct {
	Name    string
	Payload interface{}
}

const MessageProtocol = "/msg/1.0.0"

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
	host.SetStreamHandler(MessageProtocol, p.messageHandler)
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

	s, err := p.host.NewStream(ctx, peerInfo.ID, MessageProtocol)
	if err != nil {
		return err
	}

	p.messageHandler(s)
	return nil
}

func (p *P2P) GetHostData() *HostData {
	return &HostData{
		Addr: p.Addrs()[0].String(),
		Id:   p.Id(),
	}
}

func (p *P2P) messageHandler(s network.Stream) {
	go func() {
		msg := Message{
			Name:    "Hello",
			Payload: "World",
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(msg)

		if err != nil {
			fmt.Println("Error encoding message: ", err)
		}

		n, err := s.Write(buf.Bytes())
		if err != nil {
			fmt.Println("Error writing to stream: ", err)
		}
		fmt.Printf("Wrote %d bytes to stream\n", n)
		s.Close()
		/* for {
			n, err := s.Write([]byte{1})
			if err != nil {
				fmt.Println("Ping: Error writing to stream: ", err)
				return
			}

			fmt.Printf("Ping: Wrote %d bytes to stream\n", n)
		} */
	}()

	go func() {
		scanner := bufio.NewScanner(s)
		for scanner.Scan() {
			b := scanner.Bytes()
			fmt.Printf("Read %d bytes from stream: %s\n", len(b), b)
			var msg Message
			gob.NewDecoder(bytes.NewReader(b)).Decode(&msg)
			fmt.Printf("Message: %+v\n", msg)
		}

		/* 	buf := new(bytes.Buffer)
		copied, err := io.Copy(buf, s)

		if err != nil {
			fmt.Println("Error reading from stream: ", err)
		}

		fmt.Printf("Read %d bytes from stream\n", copied)

		dec := gob.NewDecoder(buf)
		var msg Message
		err = dec.Decode(&msg)

		if err != nil {
			fmt.Println("Error reading from stream: ", err)
			return
		} */

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
