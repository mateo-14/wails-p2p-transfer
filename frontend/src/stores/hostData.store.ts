import { signal } from '@preact/signals';
import { p2p } from '../../wailsjs/go/models';

export type HostDataState = (p2p.HostData & { publicAddress: string }) | null;

export const hostDataStore = signal<HostDataState | null>(null);
