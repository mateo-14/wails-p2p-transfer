import Button from '../../components/Button';
import { NavLink, Outlet } from 'react-router-dom';
import classNames from 'classnames';
import { Peer, usePeersStore } from '../../stores/peers.store';
import { useHostDataStore } from '../../stores/hostData.store';

export default function Peers() {
  const peers = usePeersStore(state => state.peers);
  const hostData = useHostDataStore(state => state.hostData);

  return (
    <div className="flex h-full">
      <div className="w-60 flex flex-col bg-zinc-900/25 flex-shrink-0">
        <div className="py-4 px-4 flex justify-between">
          <p className="text-2xl lett">Peers</p>
          <button className="font-medium hover:bg-zinc-700">+ add</button>
        </div>
        <ul className="h-full overflow-y-auto px-5">
          {peers
            .filter((peer: Peer) => peer.id !== hostData?.id)
            .map((peer: Peer) => (
              <li className="my-1" key={peer.id}>
                <NavLink
                  to={`/peers/${peer.id}`}
                  className={({ isActive }) =>
                    classNames('rounded-md px-2 py-1 block transition truncate', {
                      'bg-purple-600 text-white': isActive,
                      'hover:bg-purple-600/30 text-white/50 hover:text-white/80': !isActive
                    })
                  }
                >
                  {peer.name ?? `${peer.address}/${peer.id}`}
                </NavLink>
              </li>
            ))}
        </ul>
        <Button className="mx-2 mb-4" onClick={() => alert(hostData?.id)}>
          Your peer info
        </Button>
      </div>
      <Outlet />
    </div>
  );
}
