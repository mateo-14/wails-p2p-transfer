// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {data} from '../models';
import {main} from '../models';

export function AddFiles():Promise<Array<data.File>>;

export function ConnectToNode(arg1:string,arg2:string):Promise<void>;

export function GetPeerSharedFiles(arg1:string):Promise<Array<main.PeerFile>>;

export function OnFrontendLoad():Promise<main.InitialData>;
