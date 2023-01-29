package main

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
)

type AppNotifiee struct {
}

func (n *AppNotifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
}

func (n *AppNotifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
}

func (n *AppNotifiee) Connected(net network.Network, c network.Conn) {

	fmt.Println("Connected")
	for _, conn := range net.Conns() {
		fmt.Println(conn.RemoteMultiaddr())
	}
	for _, peer := range net.Peerstore().Peers() {
		fmt.Println(peer)
	}
}

func (n *AppNotifiee) Disconnected(net network.Network, c network.Conn) {
	fmt.Println("Disconnected")
	for _, conn := range net.Conns() {
		fmt.Println(conn.RemoteMultiaddr())
	}
	for _, peer := range net.Peerstore().Peers() {
		fmt.Println(peer)
	}
}
