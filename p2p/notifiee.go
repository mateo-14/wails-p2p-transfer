package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
)

type Notifiee struct {
}

func (n *Notifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
}

func (n *Notifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
}

func (n *Notifiee) Connected(_ network.Network, _ network.Conn) {
	fmt.Println("Connected")
}

func (n *Notifiee) Disconnected(_ network.Network, _ network.Conn) {
	fmt.Println("Disconnected")

}
