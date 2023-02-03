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
  setPeers: (peers: Peer[]) => void;
  addPeer: (peer: Peer) => void;
};

export const usePeersStore = create<PeersState>((set, get) => ({
  peers: [],

  updatePeerState: (id: string, newState: Peer['state']) =>
    set(state => ({
      peers: state.peers.map(peer => (peer.id === id ? { ...peer, state: newState } : peer))
    })),

  getPeer: (id: string) => get().peers.find(peer => peer.id === id),

  setPeers: (peers: Peer[]) => set({ peers }),
  
  addPeer: (peer: Peer) => set(state => ({ peers: [...state.peers, peer] }))
}));
