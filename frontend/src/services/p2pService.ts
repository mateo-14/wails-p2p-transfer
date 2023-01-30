import { ConnectToNode, GetPeerSharedFiles } from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime';
export function connectToPeer(address: string, id: string) {
  return ConnectToNode(address, id);
}

export type PeerConnectedCallback = (id: string) => void;
export type PeerDisconnectedCallback = (id: string) => void;

export function onPeerConnected(callback: PeerConnectedCallback) {
  EventsOn('peer:connected', callback);
}

export function onPeerDisconnected(callback: PeerDisconnectedCallback) {
  EventsOn('peer:disconnected', callback);
}

export function unsubscribePeerConnected() {
  EventsOff('peer:connected');
}

export function unsubscribePeerDisconnected() {
  EventsOff('peer:disconnected');
}

export function getPeerFiles(id: string) {
  return GetPeerSharedFiles(id);
}
