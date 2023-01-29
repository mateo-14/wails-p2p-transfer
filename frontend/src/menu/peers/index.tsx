import classNames from 'classnames';
import { PeerState, peersStore } from '../../stores/peers.store';
import { h } from 'preact';
import Router, { Link } from 'preact-router';
import Peer from './Peer';
import Button from "../../components/Button";
import { hostDataStore } from "../../stores/hostData.store";

type PeersProp = {
  path: string;
  rest?: string;
};

export default function Peers({rest}: PeersProp) {
  return (
    <div class="flex h-full">
      <div class="w-60 flex flex-col select-none bg-zinc-900/25 flex-shrink-0">
        <div class="py-4 px-4 flex justify-between">
          <p class="text-2xl lett">Peers</p>
          <button class="font-medium hover:bg-zinc-700">+ add</button>
        </div>
        <ul class="h-full overflow-y-auto px-5">
          {peersStore.value.filter(peer => peer.id !== hostDataStore.value?.id).map((peer: PeerState) => (
            <li
              class={classNames('my-1 cursor-pointer rounded-md transition truncate', {
                'bg-purple-600 text-white': peer.id === rest,
                'hover:bg-purple-600/30 text-white/50 hover:text-white/80': peer.id !== rest
              })}
              key={peer.id}
            >
              <Link href={`/peers/${peer.id}`} class="px-2 py-1 block">{peer.name ?? `${peer.address}/${peer.id}`}</Link>
            </li>
          ))}
        </ul>
        <Button className="mx-2 mb-4" onClick={() => alert(hostDataStore.value?.id)}>Your peer info</Button>
      </div>
      <Router>
        <Peer path="/peers/:id" />
      </Router>
    </div>
  );
}
