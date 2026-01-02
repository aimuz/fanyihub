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

}

