export namespace main {
	
	export class Direction {
	    id: string;
	    label: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Direction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.description = source["description"];
	    }
	}
	export class UserFeedback {
	    selectedDirections: string[];
	    userInput: string;
	
	    static createFrom(source: any = {}) {
	        return new UserFeedback(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.selectedDirections = source["selectedDirections"];
	        this.userInput = source["userInput"];
	    }
	}
	export class HistoryItem {
	    iterationId: string;
	    optimizedPrompt: string;
	    reviewReport: string;
	    evaluationReport: string;
	    score: number;
	    userFeedback: UserFeedback;
	
	    static createFrom(source: any = {}) {
	        return new HistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.iterationId = source["iterationId"];
	        this.optimizedPrompt = source["optimizedPrompt"];
	        this.reviewReport = source["reviewReport"];
	        this.evaluationReport = source["evaluationReport"];
	        this.score = source["score"];
	        this.userFeedback = this.convertValues(source["userFeedback"], UserFeedback);
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
	export class IterationData {
	    iterationId: string;
	    optimizedPrompt: string;
	    reviewReport: string;
	    evaluationReport: string;
	    score: number;
	    suggestedDirections?: Direction[];
	
	    static createFrom(source: any = {}) {
	        return new IterationData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.iterationId = source["iterationId"];
	        this.optimizedPrompt = source["optimizedPrompt"];
	        this.reviewReport = source["reviewReport"];
	        this.evaluationReport = source["evaluationReport"];
	        this.score = source["score"];
	        this.suggestedDirections = this.convertValues(source["suggestedDirections"], Direction);
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
	export class InputData {
	    version: number;
	    originalPrompt: string;
	    current: IterationData;
	    history: HistoryItem[];
	
	    static createFrom(source: any = {}) {
	        return new InputData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.originalPrompt = source["originalPrompt"];
	        this.current = this.convertValues(source["current"], IterationData);
	        this.history = this.convertValues(source["history"], HistoryItem);
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
	

}

