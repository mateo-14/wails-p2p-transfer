export namespace data {
	
	export class File {
	    path: string;
	    name: string;
	    size: number;
	    id: number;
	    hash: string;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.id = source["id"];
	        this.hash = source["hash"];
	    }
	}
	export class Peer {
	    id: number;
	    name: string;
	    address: string;
	    peerID: string;
	
	    static createFrom(source: any = {}) {
	        return new Peer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.peerID = source["peerID"];
	    }
	}
	export class InitialData {
	    hostData: p2p.HostData;
	    sharedFiles: File[];
	    peers: Peer[];
	
	    static createFrom(source: any = {}) {
	        return new InitialData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostData = this.convertValues(source["hostData"], p2p.HostData);
	        this.sharedFiles = this.convertValues(source["sharedFiles"], File);
	        this.peers = this.convertValues(source["peers"], Peer);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class PeerFile {
	    name: string;
	    size: number;
	    id: number;
	    hash: string;
	
	    static createFrom(source: any = {}) {
	        return new PeerFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
	        this.id = source["id"];
	        this.hash = source["hash"];
	    }
	}

}

export namespace p2p {
	
	export class HostData {
	    address: string;
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new HostData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.id = source["id"];
	    }
	}

}

