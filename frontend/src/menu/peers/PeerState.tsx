import classNames from 'classnames';
import type { Peer } from '../../stores/peers.store';

type PeerInfoProps = {
  peer: Peer;
};

export default function PeerInfo({ peer }: PeerInfoProps) {
  return (
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
  );
}
