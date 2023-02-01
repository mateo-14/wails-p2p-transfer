import classNames from 'classnames';
import { useEffect, useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { OnFrontendLoad } from '../wailsjs/go/main/App';
import { AppError, ErrP2PAlreadyStarted } from './errors';

import {
  onPeerConnected,
  onPeerDisconnected,
  unsubscribePeerConnected,
  unsubscribePeerDisconnected
} from './services/p2pService';
import { useHostDataStore } from './stores/hostData.store';
import { usePeersStore } from './stores/peers.store';

const IP_REGEX = '/((25[0-5]|(2[0-4]|1\\d|[1-9]|)\\d).?\\b){4}/';

export function App(props: any) {
  const [isLoading, setIsLoading] = useState(true);
  const setHostData = useHostDataStore(state => state.setHostData);
  const updatePeerState = usePeersStore(state => state.updatePeerState);

  useEffect(() => {
    OnFrontendLoad()
      .then(data => {
        fetch('https://api.ipify.org')
          .then(res => res.text())
          .then(ip => {
            const publicAddress = data.address.replace(new RegExp(IP_REGEX), `/${ip}/`);

            setHostData({
              ...data,
              address: `${data.address}/p2p`,
              publicAddress: `${publicAddress}/p2p`
            });

            onPeerConnected(id => {
              updatePeerState(id, 'connected');
            });

            onPeerDisconnected(id => {
              updatePeerState(id, 'disconnected');
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
      <div className="h-screen flex justify-center items-center bg-zinc-800">
        <p className="animate-pulse text-3xl">Loading...</p>
      </div>
    );
  }

  return (
    <div className="h-screen flex bg-zinc-800">
      {/* Menu */}
      <div className="w-16 bg-zinc-900/70 h-full flex flex-col items-center pt-4 gap-y-2">
        <MenuItem to="/peers">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="38"
            height="36"
            viewBox="0 0 38 36"
            fill="none"
          >
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M22 22.4615C22 21.8087 22.2593 21.1826 22.721 20.721C23.1826 20.2593 23.8087 20 24.4615 20H35.5385C36.1913 20 36.8174 20.2593 37.279 20.721C37.7407 21.1826 38 21.8087 38 22.4615V30.4615C38 31.1144 37.7407 31.7405 37.279 32.2021C36.8174 32.6637 36.1913 32.9231 35.5385 32.9231H33.0769V33.1339C33.0769 33.6238 33.2714 34.0939 33.6176 34.4394L34.1272 34.9497C34.2131 35.0358 34.2717 35.1454 34.2954 35.2647C34.3191 35.384 34.3069 35.5077 34.2603 35.62C34.2138 35.7324 34.135 35.8285 34.0339 35.8961C33.9328 35.9637 33.8139 35.9999 33.6923 36H26.3077C26.1861 35.9999 26.0672 35.9637 25.9661 35.8961C25.865 35.8285 25.7862 35.7324 25.7397 35.62C25.6931 35.5077 25.6809 35.384 25.7046 35.2647C25.7283 35.1454 25.7869 35.0358 25.8728 34.9497L26.3824 34.4394C26.7284 34.0934 26.9229 33.6241 26.9231 33.1348V32.9231H24.4615C23.8087 32.9231 23.1826 32.6637 22.721 32.2021C22.2593 31.7405 22 31.1144 22 30.4615V22.4615ZM23.2308 22.4615V28.6154C23.2308 28.9418 23.3604 29.2549 23.5913 29.4857C23.8221 29.7165 24.1351 29.8462 24.4615 29.8462H35.5385C35.8649 29.8462 36.1779 29.7165 36.4087 29.4857C36.6396 29.2549 36.7692 28.9418 36.7692 28.6154V22.4615C36.7692 22.1351 36.6396 21.8221 36.4087 21.5913C36.1779 21.3604 35.8649 21.2308 35.5385 21.2308H24.4615C24.1351 21.2308 23.8221 21.3604 23.5913 21.5913C23.3604 21.8221 23.2308 22.1351 23.2308 22.4615Z"
              fill="currentColor"
            />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M7.82622 17.0152C8.3701 16.9192 8.88881 17.2823 8.98478 17.8262L9.0993 18.4751C9.98745 23.508 14.0529 27.3738 19.124 28.0077C19.6721 28.0762 20.0608 28.576 19.9923 29.124C19.9238 29.6721 19.424 30.0608 18.876 29.9923C12.934 29.2495 8.1704 24.7198 7.12973 18.8227L7.01522 18.1738C6.91924 17.6299 7.28233 17.1112 7.82622 17.0152Z"
              fill="currentColor"
            />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M18.0077 6.87597C18.0762 6.32795 18.576 5.93922 19.124 6.00772C25.066 6.75047 29.8296 11.2802 30.8703 17.1773L30.9848 17.8262C31.0808 18.3701 30.7177 18.8888 30.1738 18.9848C29.6299 19.0808 29.1112 18.7177 29.0152 18.1738L28.9007 17.5249C28.0126 12.492 23.9471 8.62617 18.876 7.99228C18.3279 7.92378 17.9392 7.42399 18.0077 6.87597Z"
              fill="currentColor"
            />
            <path
              fillRule="evenodd"
              clipRule="evenodd"
              d="M0 2.46154C0 1.8087 0.25934 1.1826 0.720968 0.720968C1.1826 0.25934 1.8087 0 2.46154 0H13.5385C14.1913 0 14.8174 0.25934 15.279 0.720968C15.7407 1.1826 16 1.8087 16 2.46154V10.4615C16 11.1144 15.7407 11.7405 15.279 12.2021C14.8174 12.6637 14.1913 12.9231 13.5385 12.9231H11.0769V13.1339C11.0769 13.6238 11.2714 14.0939 11.6176 14.4394L12.1272 14.9497C12.2131 15.0358 12.2717 15.1454 12.2954 15.2647C12.3191 15.384 12.3069 15.5077 12.2603 15.62C12.2138 15.7324 12.135 15.8285 12.0339 15.8961C11.9328 15.9637 11.8139 15.9999 11.6923 16H4.30769C4.18606 15.9999 4.06719 15.9637 3.96609 15.8961C3.86499 15.8285 3.78619 15.7324 3.73966 15.62C3.69313 15.5077 3.68094 15.384 3.70464 15.2647C3.72834 15.1454 3.78687 15.0358 3.87282 14.9497L4.38236 14.4394C4.72838 14.0934 4.92286 13.6241 4.92308 13.1348V12.9231H2.46154C1.8087 12.9231 1.1826 12.6637 0.720968 12.2021C0.25934 11.7405 0 11.1144 0 10.4615V2.46154ZM1.23077 2.46154V8.61539C1.23077 8.94181 1.36044 9.25486 1.59125 9.48567C1.82207 9.71648 2.13512 9.84615 2.46154 9.84615H13.5385C13.8649 9.84615 14.1779 9.71648 14.4087 9.48567C14.6396 9.25486 14.7692 8.94181 14.7692 8.61539V2.46154C14.7692 2.13512 14.6396 1.82207 14.4087 1.59125C14.1779 1.36044 13.8649 1.23077 13.5385 1.23077H2.46154C2.13512 1.23077 1.82207 1.36044 1.59125 1.59125C1.36044 1.82207 1.23077 2.13512 1.23077 2.46154Z"
              fill="currentColor"
            />
          </svg>
        </MenuItem>
        <MenuItem to="/files">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="currentColor"
            className="w-6 h-6"
          >
            <path d="M19.5 21a3 3 0 003-3v-4.5a3 3 0 00-3-3h-15a3 3 0 00-3 3V18a3 3 0 003 3h15zM1.5 10.146V6a3 3 0 013-3h5.379a2.25 2.25 0 011.59.659l2.122 2.121c.14.141.331.22.53.22H19.5a3 3 0 013 3v1.146A4.483 4.483 0 0019.5 9h-15a4.483 4.483 0 00-3 1.146z" />
          </svg>
        </MenuItem>
      </div>

      <div className="flex-1 min-w-0">
        <Outlet />
      </div>
    </div>
  );
}

type MenuItemProps = {
  children: React.ReactNode;
  to: string;
};

function MenuItem({ children, to }: MenuItemProps) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        classNames('w-12 h-12 flex items-center justify-center p-3 rounded-xl', {
          'bg-purple-600 text-white': isActive,
          'hover:bg-purple-600/30 text-white/50 hover:text-white/80': !isActive
        })
      }
    >
      {children}
    </NavLink>
  );
}
