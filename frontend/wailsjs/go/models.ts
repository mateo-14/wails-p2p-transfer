export namespace main {
	
	export class PeerFile {
	    name: string;
	    size: number;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new PeerFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
	        this.path = source["path"];
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

