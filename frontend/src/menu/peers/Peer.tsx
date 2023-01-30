import { h } from 'preact';
import Button from '../../components/Button';
import { getPeer, updatePeerState } from '../../stores/peers.store';
import classNames from 'classnames';
import { connectToPeer, getPeerFiles } from '../../services/p2pService';

type PeerProp = {
  path: string;
  matches?: { id: string };
};

export default function Peer(props: PeerProp) {
  const peer = getPeer(props.matches?.id || '');

  if (!peer) {
    return <div>a</div>;
  }

  const handleConnect = () => {
    updatePeerState(peer.id, 'connecting');
    connectToPeer(peer.address, peer.id)
      .then(() => {
        getPeerFiles(peer.id).then((data) => {
            console.log(data)
        })
      })
      .catch(() => updatePeerState(peer.id, 'error'));
  };

  return (
    <div class="flex flex-col flex-1 min-w-0">
      <div class="flex items-center justify-between border-b-2 border-zinc-900/50 py-3 px-6 gap-x-6">
        <div class="flex flex-col min-w-0">
          <p class="text-xl font-semibold">
            {peer.name ?? peer.address}
            <span
              class={classNames('ml-4 h-2 w-2 rounded-full inline-block', {
                'animate-pulse': peer.state === 'connecting',
                'bg-green-600': peer.state === 'connected',
                'bg-cyan-500': peer.state === 'connecting',
                'bg-zinc-500': peer.state === 'disconnected' || peer.state === 'error'
              })}
            ></span>
          </p>
          <p class="text-white/50 text-sm truncate" title={`${peer.address}/${peer.id}`}>
            {peer.address}/{peer.id}
          </p>
        </div>
        {/* <div class="flex gap-x-3">
            <Button>Block</Button>
            <Button>Delete</Button>
        </div> */}
      </div>
      <div class="py-2 px-6 flex-1">
        {peer.state === 'connected' ? (
          'connected'
        ) : (
          <div class="flex items-center justify-center h-full text-3xl">
            {peer.state === 'connecting' ? <p class="animate-pulse">Connecting...</p> : null}
            {peer.state === 'disconnected' || peer.state === 'error' ? (
              <div class="flex flex-col items-center">
                <p>
                  {peer.state === 'disconnected'
                    ? 'You are disconnected from this peer'
                    : 'An error has occurred. Try again'}
                </p>
                <Button className="mt-1" onClick={handleConnect}>
                  Connect
                </Button>
              </div>
            ) : null}
          </div>
        )}
      </div>
    </div>
  );
}
