package main

import (
	"context"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func NewAppNotifiee(ctx context.Context) *AppNotifiee {
	return &AppNotifiee{
		ctx: ctx,
	}
}

type AppNotifiee struct {
	ctx context.Context
}

func (n *AppNotifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
}

func (n *AppNotifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
}

func (n *AppNotifiee) Connected(net network.Network, c network.Conn) {
	runtime.LogInfof(n.ctx, "Connected to %s\n", c.RemotePeer())
	runtime.EventsEmit(n.ctx, "peer:connected", c.RemotePeer())
}

func (n *AppNotifiee) Disconnected(net network.Network, c network.Conn) {
	runtime.LogInfof(n.ctx, "Disconnected from %s\n", c.RemotePeer())
	runtime.EventsEmit(n.ctx, "peer:disconnected", c.RemotePeer())
}
