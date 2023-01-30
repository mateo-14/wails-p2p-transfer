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

type P2P struct {
	host     host.Host
	Notifiee network.Notifiee
	msghcb   MessageHandlerCb
}

type HostData struct {
	Addr string `json:"address"`
	Id   string `json:"id"`
}

type MessageHandlerCb func(*MessageHandler)

const MessageProtocol = "/msg/1.0.0"

func NewP2P(ctx context.Context, notifiee network.Notifiee, msghcb MessageHandlerCb) (*P2P, error) {
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

	p.msghcb = msghcb

	err = p.start(ctx, privk)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *P2P) start(ctx context.Context, privk crypto.PrivKey) error {
	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/4000"), libp2p.Identity(privk))

	if p.Notifiee != nil {
		fmt.Printf("%p\n", p.Notifiee)
		host.Network().Notify(p.Notifiee)
	}

	p.host = host

	host.SetStreamHandler(MessageProtocol, func(s network.Stream) {
		msgh := NewMessageHandler(ctx, s)
		p.msghcb(msgh)
	})

	return err
}

func (p *P2P) RegenerateKey(ctx context.Context) error {
	p.host.Close()

	privk, err := generatePrivateKey()
	if err != nil {
		return err
	}

	err = storePrivateKey(privk)
	if err != nil {
		return err
	}

	return p.start(ctx, privk)
}

func (p *P2P) Addrs() []multiaddr.Multiaddr {
	return p.host.Addrs()
}

func (p *P2P) Id() string {
	return p.host.ID().String()
}

type PayloadTest struct {
	Text string
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

	payload := PayloadTest{Text: "Hello"}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(payload)

	res, _ := p.SendMessage(ctx, peerInfo.ID, Message{ID: "test", Payload: buf.Bytes()})
	fmt.Println(res.Payload)
	return nil
}

func (p *P2P) GetHostData() *HostData {
	return &HostData{
		Addr: p.Addrs()[0].String(),
		Id:   p.Id(),
	}
}

func (p *P2P) SendMessage(ctx context.Context, peerID peer.ID, msg Message) (*Message, error) {
	s, err := p.host.NewStream(ctx, peerID, MessageProtocol)

	if err != nil {
		runtime.LogErrorf(ctx, "SendMessage: Error creating stream: %s\n", err.Error())
		return nil, err
	}

	msgb, err := messageToBytes(&msg)
	if err != nil {
		runtime.LogErrorf(ctx, "SendMessage: Error encoding message: %s\n", err.Error())
		return nil, err
	}

	_, err = s.Write(msgb)
	if err != nil {
		runtime.LogErrorf(ctx, "SendMessage: Error writing to stream:%s\n ", err.Error())
		return nil, err
	}

	var res Message
	err = messageBToMessage(s, &res)
	if err != nil {
		runtime.LogErrorf(ctx, "SendMessage: Error reading response: %s\n", err.Error())
		return nil, err
	}

	return &res, nil
}

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
