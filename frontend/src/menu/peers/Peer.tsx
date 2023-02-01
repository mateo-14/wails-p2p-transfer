import classNames from 'classnames';
import { useParams } from 'react-router-dom';
import Button from '../../components/Button';
import { connectToPeer, getPeerFiles } from '../../services/p2pService';
import { usePeersStore } from '../../stores/peers.store';

export default function Peer() {
  const { id } = useParams();
  const getPeer = usePeersStore(state => state.getPeer);
  const updatePeerState = usePeersStore(state => state.updatePeerState);
  const peer = getPeer(id || '');

  if (!peer) {
    return <div>a</div>;
  }

  const handleConnect = () => {
    updatePeerState(peer.id, 'connecting');
    connectToPeer(peer.address, peer.id)
      .then(() => {
        getPeerFiles(peer.id).then(data => {
          console.log(data);
        });
      })
      .catch(() => updatePeerState(peer.id, 'error'));
  };

  return (
    <div className="flex flex-col flex-1 min-w-0">
      <div className="flex items-center justify-between border-b-2 border-zinc-900/50 py-3 px-6 gap-x-6">
        <div className="flex flex-col min-w-0">
          <p className="text-xl font-semibold">
            {peer.name ?? peer.address}
            <span
              className={classNames('ml-4 h-2 w-2 rounded-full inline-block', {
                'animate-pulse': peer.state === 'connecting',
                'bg-green-600': peer.state === 'connected',
                'bg-cyan-500': peer.state === 'connecting',
                'bg-zinc-500': peer.state === 'disconnected' || peer.state === 'error'
              })}
            ></span>
          </p>
          <p className="text-white/50 text-sm truncate" title={`${peer.address}/${peer.id}`}>
            {peer.address}/{peer.id}
          </p>
        </div>
        <div className="flex gap-x-3">
          <Button>Block</Button>
          <Button>Delete</Button>
        </div>
      </div>
      <div className="py-2 px-6 flex-1">
        {peer.state === 'connected' ? (
          'connected'
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
