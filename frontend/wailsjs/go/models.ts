export namespace channel {
	
	export class ChannelInfo {
	    id: string;
	    name: string;
	    type: string;
	    is_default: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ChannelInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.is_default = source["is_default"];
	    }
	}

}

export namespace main {
	
	export class Peer {
	    id: string;
	    username: string;
	    addr: string;
	    signaling_addr: string;
	
	    static createFrom(source: any = {}) {
	        return new Peer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.addr = source["addr"];
	        this.signaling_addr = source["signaling_addr"];
	    }
	}

}

export namespace config {
	
	export class TURNConfig {
	    URL: string;
	    Username: string;
	    Password: string;
	
	    static createFrom(source: any = {}) {
	        return new TURNConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.URL = source["URL"];
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	    }
	}
	export class Config {
	    NetworkMode: string;
	    AppMode: string;
	    Port: number;
	    DataDir: string;
	    STUNURLs: string[];
	    TURNConfig: TURNConfig;
	    ServerAddr: string;
	    Username: string;
	    StorageType: string;
	    MongoDBURI: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.NetworkMode = source["NetworkMode"];
	        this.AppMode = source["AppMode"];
	        this.Port = source["Port"];
	        this.DataDir = source["DataDir"];
	        this.STUNURLs = source["STUNURLs"];
	        this.TURNConfig = this.convertValues(source["TURNConfig"], TURNConfig);
	        this.ServerAddr = source["ServerAddr"];
	        this.Username = source["Username"];
	        this.StorageType = source["StorageType"];
	        this.MongoDBURI = source["MongoDBURI"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

}

