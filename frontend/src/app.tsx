import { StartP2P, ConnectToNode } from '../wailsjs/go/main/App';
import { useState } from 'preact/hooks';
import { h } from 'preact';
import { main } from '../wailsjs/go/models';
import { JSXInternal } from 'preact/src/jsx';

type HostDataState = (main.HostData & { publicAddress: string }) | null;

const IP_REGEX = '/((25[0-5]|(2[0-4]|1\\d|[1-9]|)\\d).?\\b){4}/';
export function App(props: any) {
  const [hostData, setHostData] = useState<HostDataState>(null);

  const startP2P = () => {
    setHostData(null);

    StartP2P()
      .then(data => {
        fetch('https://api.ipify.org')
          .then(res => res.text())
          .then(ip => {
            const publicAddress = data.address.replace(new RegExp(IP_REGEX), `/${ip}/`);

            setHostData({
              ...data,
              address: `${data.address}/p2p/${data.id}}`,
              publicAddress: `${publicAddress}/p2p/${data.id}`,
            });
          });
      })
      .catch(err => {});
  };

  const copyPublicAddress = () => {
    if (hostData === null) return;

    navigator.clipboard.writeText(`${hostData.publicAddress}/${hostData.id}`);
  };


  const submitConnect = (e: h.JSX.TargetedEvent<HTMLFormElement, Event>) => {
    e.preventDefault();
    const address = e.currentTarget.address.value;
    if (address === '') return;

    ConnectToNode(address)
      .then(() => {
        console.log('Connected');
      })
      .catch(err => {
        console.log(err);
      });
  };
  return (
    <div class="min-h-screen bg-zinc-900">
      {hostData ? (
        <>
          <div>
              <p>Your node ID is: {hostData.id}</p>
              <p>
                Addresses: local: {hostData.address}, public: {hostData.publicAddress}
              </p>
              <p>

              </p>
          </div>
          <Button onClick={copyPublicAddress}>Copy public address to share</Button>

          <form onSubmit={submitConnect} class="w-full">
            <label htmlFor="address-input">Connect to: </label>
            <input type="text" id="address-input" class="text-black w-full" name="address" />
            <Button>Connect</Button>
          </form>
        </>
      ) : (
        <Button onClick={startP2P}>Start P2P</Button>
      )}
    </div>
  );
}

function Button(props: JSXInternal.IntrinsicElements['button']) {
  return (
    <button
      class="bg-purple-700 py-2 px-3 rounded-md text-sm font-semibold hover:bg-purple-600 hover:shadow-lg active:shadow-lg hover:shadow-purple-600/20 active:shadow-purple-600/50 transition-all"
      {...props}
    >
      {props.children}
    </button>
  );
}
