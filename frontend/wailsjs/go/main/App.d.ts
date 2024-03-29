// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {data} from '../models';

export function AddFiles():Promise<string>;

export function AddPeer(arg1:string,arg2:string):Promise<data.Peer>;

export function ConnectToNode(arg1:string,arg2:string):Promise<void>;

export function DownloadFile(arg1:string,arg2:number):Promise<void>;

export function GetPeerSharedFiles(arg1:string):Promise<Array<data.PeerFile>>;

export function OnFrontendLoad():Promise<data.InitialData>;

export function RemoveSharedFile(arg1:number):Promise<void>;
