import { create } from 'zustand';
import { p2p } from '../../wailsjs/go/models';

export type HostData = p2p.HostData & { publicAddress: string };
export type HostDataState = {
  hostData: HostData | null;
  setHostData: (hostData: HostData) => void;
};

export const useHostDataStore = create<HostDataState>((set, get) => ({
  hostData: null,
  setHostData: (hostData: HostData) => set({ hostData })
}));
