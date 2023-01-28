export namespace main {
	
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

