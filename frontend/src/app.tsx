import { OnFrontendLoad } from '../wailsjs/go/main/App';
import { useEffect, useState } from 'preact/hooks';
import { h } from 'preact';
import { AppError, ErrP2PAlreadyStarted } from './errors';
import { Router } from 'preact-router';
import Peers from './menu/peers';
import { hostDataStore } from './stores/hostData.store';
import {
  onPeerConnected,
  onPeerDisconnected,
  unsubscribePeerConnected,
  unsubscribePeerDisconnected
} from './services/p2pService';
import { updatePeerState } from "./stores/peers.store";

const IP_REGEX = '/((25[0-5]|(2[0-4]|1\\d|[1-9]|)\\d).?\\b){4}/';
export function App(props: any) {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    OnFrontendLoad()
      .then(data => {
        fetch('https://api.ipify.org')
          .then(res => res.text())
          .then(ip => {
            const publicAddress = data.address.replace(new RegExp(IP_REGEX), `/${ip}/`);
            hostDataStore.value = {
              ...data,
              address: `${data.address}/p2p`,
              publicAddress: `${publicAddress}/p2p`
            };

            onPeerConnected(id => {
              updatePeerState(id, 'connected')
            });

            onPeerDisconnected(id => {
              updatePeerState(id, 'disconnected')
            });

            setIsLoading(false);
          });
      })
      .catch((err: AppError) => {
        console.log(err);
        if (err === ErrP2PAlreadyStarted) {
          setIsLoading(false);
        }
      });

    return () => {
      unsubscribePeerConnected();
      unsubscribePeerDisconnected();
    };
  }, []);

  if (isLoading) {
    return (
      <div class="h-screen flex justify-center items-center bg-zinc-800">
        <p class="animate-pulse text-3xl">Loading...</p>
      </div>
    );
  }

  return (
    <div class="h-screen flex bg-zinc-800">
      {/* Menu */}
      <div class="w-16 bg-zinc-900/70 h-full">Menu</div>

      <div class="flex-1 min-w-0">
        <Router url="/peers">
          <Peers path="/peers/:rest*" />
        </Router>
      </div>
    </div>
  );
}
