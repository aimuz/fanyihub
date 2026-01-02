export namespace types {
	
	export class DetectResult {
	    code: string;
	    name: string;
	    defaultTarget: string;
	
	    static createFrom(source: any = {}) {
	        return new DetectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.name = source["name"];
	        this.defaultTarget = source["defaultTarget"];
	    }
	}
	export class Provider {
	    name: string;
	    type: string;
	    base_url?: string;
	    api_key: string;
	    model: string;
	    system_prompt?: string;
	    max_tokens?: number;
	    temperature?: number;
	    active: boolean;
	    disable_thinking?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Provider(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.base_url = source["base_url"];
	        this.api_key = source["api_key"];
	        this.model = source["model"];
	        this.system_prompt = source["system_prompt"];
	        this.max_tokens = source["max_tokens"];
	        this.temperature = source["temperature"];
	        this.active = source["active"];
	        this.disable_thinking = source["disable_thinking"];
	    }
	}
	export class TranslateRequest {
	    text: string;
	    sourceLang: string;
	    targetLang: string;
	
	    static createFrom(source: any = {}) {
	        return new TranslateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.sourceLang = source["sourceLang"];
	        this.targetLang = source["targetLang"];
	    }
	}
	export class Usage {
	    promptTokens: number;
	    completionTokens: number;
	    totalTokens: number;
	    cacheHit: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Usage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.promptTokens = source["promptTokens"];
	        this.completionTokens = source["completionTokens"];
	        this.totalTokens = source["totalTokens"];
	        this.cacheHit = source["cacheHit"];
	    }
	}
	export class TranslateResult {
	    text: string;
	    usage: Usage;
	
	    static createFrom(source: any = {}) {
	        return new TranslateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.usage = this.convertValues(source["usage"], Usage);
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

