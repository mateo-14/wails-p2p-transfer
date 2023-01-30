package p2p

import (
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
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageType int64

const (
	Request MessageType = iota
	Response
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
	Type    MessageType
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

	p.SendMessage(ctx, peerInfo.ID, Message{Name: "Hello", Payload: "World"})
	return nil
}

func (p *P2P) GetHostData() *HostData {
	return &HostData{
		Addr: p.Addrs()[0].String(),
		Id:   p.Id(),
	}
}

func (p *P2P) messageToBytes(msg *Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(msg)
	return buf.Bytes(), err
}

func (p *P2P) SendMessage(ctx context.Context, peerID peer.ID, msg Message) error {
	s, err := p.host.NewStream(ctx, peerID, MessageProtocol)

	if err != nil {
		runtime.LogErrorf(ctx, "SendMessage: Error creating stream: %s\n", err.Error())
		return err
	}

	msgb, err := p.messageToBytes(&msg)
	if err != nil {
		runtime.LogErrorf(ctx, "SendMessage: Error encoding message: %s\n", err.Error())
		return err
	}

	n, err := s.Write(msgb)
	if err != nil {
		runtime.LogErrorf(ctx, "SendMessage: Error writing to stream:%s\n ", err.Error())
	}

	fmt.Printf("Wrote %d bytes to stream\n", n)
	p.messageHandler(s)

	return err
}

func (p *P2P) messageHandler(s network.Stream) {
	// Manage request type messages
	go func() {
		var msg Message
		dec := gob.NewDecoder(s)
		err := dec.Decode(&msg)
		if err != nil {
			fmt.Println("Error reading from stream: ", err)
			return
		}

		switch msg.Type {
		case Request:
			fmt.Printf("Request: %+v\n", msg)
			msgb, err := p.messageToBytes(&Message{
				Type:    Response,
				Name:    "Hello",
				Payload: "World",
			})

			if err != nil {
				runtime.LogErrorf(context.Background(), "SendMessage: Error encoding message: %s\n", err.Error())
				return
			}

			s.Write(msgb)
			s.Close()
		case Response:
			fmt.Printf("Response: %+v\n", msg)
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
