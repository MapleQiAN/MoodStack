export namespace app {
	
	export class Session {
	    id: string;
	    userId: number;
	    sessionKey: string;
	    // Go type: time
	    expiresAt: any;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.userId = source["userId"];
	        this.sessionKey = source["sessionKey"];
	        this.expiresAt = this.convertValues(source["expiresAt"], null);
	        this.createdAt = this.convertValues(source["createdAt"], null);
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
	export class User {
	    id: number;
	    username: string;
	    biometricEnabled: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.biometricEnabled = source["biometricEnabled"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class AuthResult {
	    success: boolean;
	    user?: User;
	    session?: Session;
	    message: string;
	    requireSetup: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AuthResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.user = this.convertValues(source["user"], User);
	        this.session = this.convertValues(source["session"], Session);
	        this.message = source["message"];
	        this.requireSetup = source["requireSetup"];
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
	export class BiometricInfo {
	    supported: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new BiometricInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.supported = source["supported"];
	        this.message = source["message"];
	    }
	}
	export class BiometricTestResult {
	    com_initialization: boolean;
	    factory_access: boolean;
	    api_available: boolean;
	    error_message: string;
	
	    static createFrom(source: any = {}) {
	        return new BiometricTestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.com_initialization = source["com_initialization"];
	        this.factory_access = source["factory_access"];
	        this.api_available = source["api_available"];
	        this.error_message = source["error_message"];
	    }
	}
	export class Diary {
	    id: string;
	    title: string;
	    content: string;
	    fileName: string;
	    fileType: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new Diary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.fileName = source["fileName"];
	        this.fileType = source["fileType"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.tags = source["tags"];
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
	export class DiaryEncryptionInfo {
	    mode: string;
	    hasSalt: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DiaryEncryptionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.hasSalt = source["hasSalt"];
	    }
	}
	export class DiaryEncryptionOptions {
	    mode: string;
	    individualPassword?: string;
	
	    static createFrom(source: any = {}) {
	        return new DiaryEncryptionOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.individualPassword = source["individualPassword"];
	    }
	}
	export class EncryptedDiary {
	    id: string;
	    userId: number;
	    title: string;
	    fileName: string;
	    fileType: string;
	    encryptionMode: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new EncryptedDiary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.userId = source["userId"];
	        this.title = source["title"];
	        this.fileName = source["fileName"];
	        this.fileType = source["fileType"];
	        this.encryptionMode = source["encryptionMode"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class MigrationStatus {
	    migrationNeeded: boolean;
	    diaryCount: number;
	
	    static createFrom(source: any = {}) {
	        return new MigrationStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.migrationNeeded = source["migrationNeeded"];
	        this.diaryCount = source["diaryCount"];
	    }
	}
	export class SearchResult {
	    diary: Diary;
	    matchedSnippets: string[];
	    matchType: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.diary = this.convertValues(source["diary"], Diary);
	        this.matchedSnippets = source["matchedSnippets"];
	        this.matchType = source["matchType"];
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

export namespace main {
	
	export class AuthStatusResult {
	    isAuthenticated: boolean;
	    hasUsers: boolean;
	    requireSetup: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AuthStatusResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isAuthenticated = source["isAuthenticated"];
	        this.hasUsers = source["hasUsers"];
	        this.requireSetup = source["requireSetup"];
	    }
	}

}

