import { create } from 'zustand';

export type Peer = {
  id: string;
  address: string;
  name: string;
  state: 'connected' | 'disconnected' | 'connecting' | 'error';
};

export type PeersState = {
  peers: Peer[];
  updatePeerState: (id: string, newState: Peer['state']) => void;
};

export const usePeersStore = create<PeersState>((set, get) => ({
  peers: [
    {
      id: 'QmWXNcTtpYnHKSVEcUPRCvMs3yx7caqGqBZrEB9hRghj6z',
      address: '/ip4/192.168.0.176/tcp/4000/p2p',
      name: 'Desktop',
      state: 'disconnected'
    },
    {
      id: 'QmdberyhsY3DGE2pVjY1Etw7q3xiToHLb1BvXFvHySdACD',
      address: '/ip4/192.168.0.247/tcp/4000/p2p',
      name: 'Laptop',
      state: 'error'
    }
  ],

  updatePeerState: (id: string, newState: Peer['state']) =>
    set(state => ({
      peers: state.peers.map(peer => (peer.id === id ? { ...peer, state: newState } : peer))
    })),

  getPeer: (id: string) => get().peers.find(peer => peer.id === id)
}));
