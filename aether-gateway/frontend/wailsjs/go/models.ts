export namespace main {
	
	export class GatewaySettings {
	    password_hash: string;
	    disk_limit_gb: number;
	    theme: string;
	    mode: string;
	
	    static createFrom(source: any = {}) {
	        return new GatewaySettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.password_hash = source["password_hash"];
	        this.disk_limit_gb = source["disk_limit_gb"];
	        this.theme = source["theme"];
	        this.mode = source["mode"];
	    }
	}

}

