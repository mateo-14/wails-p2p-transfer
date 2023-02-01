import { ConnectToNode, GetPeerSharedFiles } from '../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../wailsjs/runtime';

const PeerDisconnectedEvent = 'peer:disconnected';
const PeerConnectedEvent = 'peer:connected';

export function connectToPeer(address: string, id: string) {
  return ConnectToNode(address, id);
}

export type PeerConnectedCallback = (id: string) => void;
export type PeerDisconnectedCallback = (id: string) => void;

export function onPeerConnected(callback: PeerConnectedCallback) {
  EventsOn(PeerConnectedEvent, callback);
}

export function onPeerDisconnected(callback: PeerDisconnectedCallback) {
  EventsOn(PeerDisconnectedEvent, callback);
}

export function unsubscribePeerConnected() {
  EventsOff(PeerConnectedEvent);
}

export function unsubscribePeerDisconnected() {
  EventsOff(PeerDisconnectedEvent);
}

export function getPeerFiles(id: string) {
  return GetPeerSharedFiles(id);
}
