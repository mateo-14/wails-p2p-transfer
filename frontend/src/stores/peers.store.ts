import { signal } from '@preact/signals';

export type PeerState = {
  id: string;
  address: string;
  name: string;
  state: 'connected' | 'disconnected' | 'connecting' | 'error';
};

export const peersStore = signal<PeerState[]>([
  {
    id: 'QmWXNcTtpYnHKSVEcUPRCvMs3yx7caqGqBZrEB9hRghj6z',
    address: '/ip4/192.168.0.176/tcp/4000/p2p',
    name: 'Desktop',
    state: 'disconnected'
  },
  {
    id: 'QmQzCe3gFukuvLaZao1ZPmffUbMvvRMjwa9vguZB8uLkWU',
    address: '/ip4/192.168.0.247/tcp/4000/p2p',
    name: 'Laptop',
    state: 'error'
  }
]);

export function updatePeerState(id: string, state: PeerState['state']) {
  peersStore.value = peersStore.value.map(peer => {
    if (peer.id === id) {
      return {
        ...peer,
        state
      };
    }
    return peer;
  });
}

export function getPeer(id: string) {
  return peersStore.value.find(peer => peer.id === id);
}