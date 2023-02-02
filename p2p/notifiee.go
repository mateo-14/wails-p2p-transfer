package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func NewNotifiee(ctx context.Context, host host.Host) *Notifiee {
	return &Notifiee{
		ctx:  ctx,
		Host: host,
	}
}

type Notifiee struct {
	ctx context.Context
	host.Host
}

func (n *Notifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {

}

func (n *Notifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
}

func (n *Notifiee) Connected(net network.Network, c network.Conn) {
	net.Peerstore().AddAddr(c.RemotePeer(), c.RemoteMultiaddr(), peerstore.RecentlyConnectedAddrTTL)
	runtime.LogInfof(n.ctx, "Connected to %s\n", c.RemotePeer())
}

func (n *Notifiee) Disconnected(net network.Network, c network.Conn) {
	runtime.LogInfof(n.ctx, "Disconnected from %s\n", c.RemotePeer())
}
