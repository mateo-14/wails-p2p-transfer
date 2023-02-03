import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import Button from '../../components/Button';
import { connectToPeer, getPeerFiles } from '../../services/p2pService';
import { usePeersStore } from '../../stores/peers.store';
import PeerFiles from './PeerFiles';
import PeerInfo from "./PeerState";
import { data } from '../../../wailsjs/go/models';

export default function Peer() {
  const { id } = useParams();
  const peer = usePeersStore(state => state.peers.find(peer => peer.id === id));
  const updatePeerState = usePeersStore(state => state.updatePeerState);
  const [files, setFiles] = useState<data.PeerFile[]>([]);

  if (!id || !peer) {
    return <div>Not found</div>;
  }

  const handleConnect = () => {
    updatePeerState(peer.id, 'connecting');
    connectToPeer(peer.address, peer.id).catch(() => updatePeerState(peer.id, 'error'));
  };

  useEffect(() => {
    if (peer.state === 'connected') {
      getPeerFiles(peer.id).then(data => {
        if (data) setFiles(data);
      });
    }
  }, [peer.state]);

  return (
    <div className="flex flex-col flex-1 min-w-0">
      <div className="flex items-center justify-between border-b-2 border-zinc-900/50 py-3 px-6 gap-x-6">
        <PeerInfo peer={peer} />
      </div>
      <div className="py-2 pl-6 flex-1 min-h-0">
        {peer.state === 'connected' ? (
          <PeerFiles files={files} peerId={id} />
        ) : (
          <div className="flex items-center justify-center h-full text-3xl">
            {peer.state === 'connecting' ? <p className="animate-pulse">Connecting...</p> : null}
            {peer.state === 'disconnected' || peer.state === 'error' ? (
              <div className="flex flex-col items-center">
                <p>
                  {peer.state === 'disconnected'
                    ? 'You are disconnected from this peer'
                    : 'An error has occurred. Try again'}
                </p>
                <Button className="mt-1" onClick={handleConnect}>
                  Connectd
                </Button>
              </div>
            ) : null}
          </div>
        )}
      </div>
    </div>
  );
}
