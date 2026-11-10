import { apiConfig } from "../config/envConfig.js";
import type { AppState, AuthState, StateListener, User } from "./types";
import {
    AUTH_ALIAS_KEYS,
    getAuthAlias,
    normalizeAuth,
    normalizePermissions,
    normalizeRoles,
    syncDerivedAuthState,
    updateAuthUserProperty
} from "./auth";
import { createListenerRegistry } from "./listeners";
import {
    clearRouteState,
    deleteRouteState,
    getRouteCache,
    getRouteStateMap,
    migrateLegacyToken,
    readLocalStorage,
    readPersistentJSON,
    readPersistentNumber,
    readRouteState,
    readSessionStorage,
    readStorage,
    removeStorage,
    restoreScroll,
    routeCache,
    routeState,
    saveScroll,
    writeLocalStorage,
    writeSessionStorage,
    writeStorage,
    writeRouteState
} from "./storage";

export const {
    MAIN_URL,
    EMBED_URL,
    BANNERDROP_URL,
    API_URL,
    STRIPE_URL,
    AD_URL,
    SEARCH_URL,
    MERE_URL,
    MERE_WS,
    CHAT_URL,
    CHAT_WS,
    MUSIC_URL,
    LIVE_URL,
    SRC_URL,
    FILEDROP_URL,
    CHATDROP_URL
} = apiConfig;

const allowedKeys = new Set<string>([
    "token", "user", "username", "userProfile", "socket",
    "favFarms",
    "roles", "permissions", "auth", "environment", "lang",
    "lastPath", "currentRoute", "routeCache", "routeState",
    "currentChatId", "isLoading", "userid", "unreadMessages",
    "unreadNotifications", "isLoggedIn"
]);

const PERSISTED_KEYS = new Set<string>([
    "userProfile", "user", "roles", "permissions",
    "favFarms",
    "unreadMessages", "unreadNotifications"
]);

const SESSION_KEYS = new Set<string>(["token"]);

const ROUTE_CACHE_KEY = "routeCache";
const ROUTE_STATE_KEY = "routeState";

const stateListenerRegistry = createListenerRegistry<AppState>(
    () => state,
    (key, currentState) => {
        if (AUTH_ALIAS_KEYS.has(key)) {
            return getAuthAlias(currentState, key);
        }
        if (key === ROUTE_CACHE_KEY) {
            return routeCache;
        }
        if (key === ROUTE_STATE_KEY) {
            return routeState;
        }
        return currentState[key];
    }
);
const scheduleNotify = stateListenerRegistry.scheduleNotify;

function getValueByPath(path: string, source: any = state): any {
    if (!path) {
        return source;
    }
    return path.split(".").reduce((current, part) => current?.[part], source);
}

const initialToken = readSessionStorage("token");
const initialUser = readPersistentJSON<User>("user", null);
const initialProfile = readPersistentJSON<Record<string, any>>("userProfile", {});
const initialRoles = normalizeRoles(readPersistentJSON("roles", []));
const initialPermissions = normalizePermissions(readPersistentJSON("permissions", []));
const initialUnreadMessages = readPersistentNumber("unreadMessages", 0);
const initialUnreadNotifications = readPersistentNumber("unreadNotifications", 0);
const initialFavFarms = readPersistentJSON<string[]>("favFarms", []);

const initialAuth = normalizeAuth({
    accessToken: initialToken || null,
    user: initialUser || null,
    roles: initialRoles,
    permissions: initialPermissions,
    isAuthenticated: Boolean(initialToken),
    loading: false
});

const rawState: AppState = {
    auth: initialAuth,
    userProfile: initialProfile && typeof initialProfile === "object" ? initialProfile : {},
    socket: null,
    environment: {},
    lang: "en",
    lastPath: typeof window !== "undefined" ? window.location.pathname + window.location.search : "/",
    currentRoute: null,
    routeCache,
    routeState,
    currentChatId: null,
    isLoading: false,
    unreadMessages: initialUnreadMessages,
    unreadNotifications: initialUnreadNotifications,
    favFarms: Array.isArray(initialFavFarms) ? initialFavFarms : [],
    isLoggedIn: Boolean(initialAuth.isAuthenticated || initialUser?.id || initialUser?.userid)
};

function triggerAuthNotifications(): void {
    scheduleNotify("token");
    scheduleNotify("user");
    scheduleNotify("username");
    scheduleNotify("userid");
    scheduleNotify("roles");
    scheduleNotify("permissions");
    scheduleNotify("isLoggedIn");
    scheduleNotify("auth");
}

function setAuthAlias(key: string, value: any): void {
    const currentAuth = rawState.auth;
    switch (key) {
        case "token": {
            const token = value || null;
            rawState.auth = normalizeAuth({
                ...currentAuth,
                accessToken: token,
                isAuthenticated: Boolean(token || (currentAuth.user && typeof currentAuth.user === "object" && (currentAuth.user.id || currentAuth.user.userid)))
            }, currentAuth);
            syncDerivedAuthState(rawState, rawState.auth);
            triggerAuthNotifications();
            break;
        }
        case "user": {
            rawState.auth = normalizeAuth({
                ...currentAuth,
                user: value || null,
                isAuthenticated: Boolean(value || currentAuth.accessToken)
            }, currentAuth);
            syncDerivedAuthState(rawState, rawState.auth);
            triggerAuthNotifications();
            break;
        }
        case "isLoggedIn": {
            const isLoggedIn = Boolean(value);
            if (!isLoggedIn) {
                rawState.auth = normalizeAuth({
                    ...currentAuth,
                    accessToken: null,
                    user: null,
                    isAuthenticated: false
                }, currentAuth);
            } else {
                rawState.auth = normalizeAuth({
                    ...currentAuth,
                    isAuthenticated: true
                }, currentAuth);
            }
            syncDerivedAuthState(rawState, rawState.auth);
            triggerAuthNotifications();
            break;
        }
        case "roles": {
            rawState.auth = normalizeAuth({
                ...currentAuth,
                roles: normalizeRoles(value)
            }, currentAuth);
            scheduleNotify("roles");
            scheduleNotify("isLoggedIn");
            scheduleNotify("auth");
            break;
        }
        case "permissions": {
            rawState.auth = normalizeAuth({
                ...currentAuth,
                permissions: normalizePermissions(value)
            }, currentAuth);
            scheduleNotify("permissions");
            scheduleNotify("auth");
            break;
        }
        case "username": {
            updateAuthUserProperty(rawState, "username", value ?? "");
            scheduleNotify("username");
            scheduleNotify("user");
            scheduleNotify("auth.user");
            scheduleNotify("isLoggedIn");
            scheduleNotify("auth");
            break;
        }
        case "userid": {
            updateAuthUserProperty(rawState, "userid", value ?? null);
            scheduleNotify("userid");
            scheduleNotify("user");
            scheduleNotify("auth.user");
            scheduleNotify("isLoggedIn");
            scheduleNotify("auth");
            break;
        }
        default:
            break;
    }
}

const proxyCache = new WeakMap<object, Map<string, any>>();

function getCachedProxy<T extends object>(target: T, path: string[]): T {
    let pathMap = proxyCache.get(target);
    if (!pathMap) {
        pathMap = new Map();
        proxyCache.set(target, pathMap);
    }
    const pathKey = path.join(".");
    const existing = pathMap.get(pathKey);
    if (existing) {
        return existing;
    }
    const proxy = createReactiveObject(target, path);
    pathMap.set(pathKey, proxy);
    return proxy;
}

function isObjectLike(value: any): boolean {
    return (value !== null && typeof value === "object");
}

function shouldProxy(value: any): boolean {
    return (isObjectLike(value) && !(value instanceof Map) && !(value instanceof Set) && !(value instanceof Date) && !(value instanceof RegExp));
}

function createReactiveObject<T extends object>(obj: T, path: string[] = []): T {
    if (!shouldProxy(obj)) {
        return obj;
    }

    return new Proxy(obj, {
        get(target, prop, receiver) {
            const key = String(prop);

            if (path.length === 0 && AUTH_ALIAS_KEYS.has(key)) {
                return getAuthAlias(rawState, key);
            }

            const value = Reflect.get(target, prop, receiver);

            if (shouldProxy(value as any)) {
                return getCachedProxy(value as any, path.concat(String(prop)));
            }

            return value;
        },
        set(target, prop, value, receiver) {
            const key = String(prop);
            if (path.length === 0 && AUTH_ALIAS_KEYS.has(key)) {
                setAuthAlias(key, value);
                return true;
            }
            if (path.length === 0 && key === "auth") {
                const previous = (target as any).auth;
                (target as any).auth = normalizeAuth(value, previous);
                triggerAuthNotifications();
                return true;
            }
            if (path.length === 0 && key === "isLoading") {
                const changed = (target as any)[key] !== Boolean(value);
                if (!changed) {
                    return true;
                }
                Reflect.set(target, prop, Boolean(value), receiver);
                (target as any).auth = normalizeAuth({
                    ...(target as any).auth,
                    loading: Boolean(value)
                }, (target as any).auth);
                scheduleNotify("isLoading");
                scheduleNotify("auth");
                return true;
            }
            const oldValue = Reflect.get(target, prop, receiver);
            if (Object.is(oldValue, value)) {
                return true;
            }
            const result = Reflect.set(target, prop, value, receiver);
            if (!result) {
                return false;
            }
            const fullPath = path.concat(key).join(".");
            scheduleNotify(fullPath);
            if (path.length > 0) {
                scheduleNotify(path[0]);
            }
            if (path[0] === "auth") {
                if (path[1] === "user") {
                    scheduleNotify("user");
                    scheduleNotify("username");
                    scheduleNotify("userid");
                    scheduleNotify("isLoggedIn");
                }
                if (path[1] === "accessToken") {
                    scheduleNotify("token");
                    scheduleNotify("isLoggedIn");
                }
                if (path[1] === "roles") {
                    scheduleNotify("roles");
                }
                if (path[1] === "permissions") {
                    scheduleNotify("permissions");
                }
            }
            return true;
        },
        deleteProperty(target, prop) {
            const key = String(prop);
            if (path.length === 0 && AUTH_ALIAS_KEYS.has(key)) {
                setAuthAlias(key, null);
                return true;
            }
            if (!Object.prototype.hasOwnProperty.call(target, prop)) {
                return true;
            }
            const deleted = Reflect.deleteProperty(target, prop);
            if (!deleted) {
                return false;
            }
            const fullPath = path.concat(key).join(".");
            scheduleNotify(fullPath);
            if (path.length > 0) {
                scheduleNotify(path[0]);
            }
            return true;
        }
    });
}

const state = getCachedProxy(rawState, []) as AppState;

function getStateValue(key: string): any {
    if (AUTH_ALIAS_KEYS.has(key)) {
        return getAuthAlias(rawState, key);
    }
    if (key === ROUTE_CACHE_KEY) {
        return routeCache;
    }
    if (key === ROUTE_STATE_KEY) {
        return routeState;
    }
    return state[key];
}

function persistStateKey(key: string, value: any): void {
    if (SESSION_KEYS.has(key)) {
        writeSessionStorage(key, value);
        return;
    }
    if (PERSISTED_KEYS.has(key)) {
        writeLocalStorage(key, value);
    }
}

function broadcastAuthChange(reason = "updated"): void {
    try {
        localStorage.setItem("auth:changed", JSON.stringify({
            reason,
            ts: Date.now()
        }));
        localStorage.removeItem("auth:changed");
    } catch {
        // Ignore
    }
}

function setAuthState(value: Partial<AuthState>, persist = false): void {
    const previous = rawState.auth;
    const next = normalizeAuth(value, previous);
    rawState.auth = next;
    syncDerivedAuthState(rawState, next);
    triggerAuthNotifications();
    if (persist) {
        persistStateKey("token", next.accessToken);
        persistStateKey("user", next.user);
        persistStateKey("roles", next.roles);
        persistStateKey("permissions", next.permissions);
        broadcastAuthChange(next.isAuthenticated ? "login" : "logout");
    }
}

function setState(keyOrObject: string | Record<string, any>, persistOrValue: boolean | any = false, maybeValue: any = undefined): void {
    const updates: Record<string, any> = typeof keyOrObject === "object" && keyOrObject !== null ? keyOrObject : {
        [keyOrObject as string]: persistOrValue
    };
    const persist = typeof keyOrObject === "object" && keyOrObject !== null ? Boolean(persistOrValue) : Boolean(maybeValue);

    const authUpdates: Partial<AuthState> = {};
    let hasAuthUpdate = false;

    for (const [key, value] of Object.entries(updates)) {
        if (!allowedKeys.has(key)) {
            throw new Error(`Invalid state key:${key}`);
        }
        if (key === ROUTE_CACHE_KEY || key === ROUTE_STATE_KEY) {
            throw new Error(`${key} is read-only. Use dedicated route cache helpers.`);
        }
        if (AUTH_ALIAS_KEYS.has(key)) {
            hasAuthUpdate = true;
            switch (key) {
                case "token":
                    authUpdates.accessToken = value || null;
                    break;
                case "user":
                    authUpdates.user = value || null;
                    break;
                case "isLoggedIn":
                    if (!value) {
                        authUpdates.accessToken = null;
                        authUpdates.user = null;
                        authUpdates.isAuthenticated = false;
                    } else {
                        authUpdates.isAuthenticated = true;
                    }
                    break;
                case "roles":
                    authUpdates.roles = normalizeRoles(value);
                    break;
                case "permissions":
                    authUpdates.permissions = normalizePermissions(value);
                    break;
                case "username":
                    if (authUpdates.user === undefined) {
                        authUpdates.user = getAuthAlias(rawState, "user");
                    }
                    if (authUpdates.user && typeof authUpdates.user === "object" && !Array.isArray(authUpdates.user)) {
                        authUpdates.user = {
                            ...authUpdates.user,
                            username: value ?? ""
                        };
                    } else {
                        authUpdates.user = {
                            username: value ?? ""
                        };
                    }
                    break;
                case "userid":
                    if (authUpdates.user === undefined) {
                        authUpdates.user = getAuthAlias(rawState, "user");
                    }
                    if (authUpdates.user && typeof authUpdates.user === "object" && !Array.isArray(authUpdates.user)) {
                        authUpdates.user = {
                            ...authUpdates.user,
                            userid: value ?? null
                        };
                    } else {
                        authUpdates.user = {
                            userid: value ?? null
                        };
                    }
                    break;
                default:
                    break;
            }
            continue;
        }
        if (key === "auth") {
            hasAuthUpdate = true;
            if (value && typeof value === "object") {
                Object.assign(authUpdates, value);
            }
            continue;
        }
        state[key] = value;
        if (persist && PERSISTED_KEYS.has(key)) {
            persistStateKey(key, value);
        }
    }

    if (hasAuthUpdate) {
        setAuthState(authUpdates, persist);
    }

    if (Object.prototype.hasOwnProperty.call(updates, "isLoading")) {
        rawState.auth = normalizeAuth({
            ...rawState.auth,
            loading: Boolean(rawState.isLoading)
        }, rawState.auth);
        scheduleNotify("auth");
    }
}

function buildPublicSnapshot(): Partial<AppState> & Record<string, any> {
    return {
        token: getAuthAlias(rawState, "token"),
        user: getAuthAlias(rawState, "user"),
        username: getAuthAlias(rawState, "username"),
        userProfile: state.userProfile,
        socket: state.socket,
        roles: getAuthAlias(rawState, "roles"),
        permissions: getAuthAlias(rawState, "permissions"),
        auth: state.auth,
        environment: state.environment,
        lang: state.lang,
        lastPath: state.lastPath,
        currentRoute: state.currentRoute,
        routeCache,
        routeState,
        currentChatId: state.currentChatId,
        isLoading: state.isLoading,
        userid: getAuthAlias(rawState, "userid"),
        unreadMessages: state.unreadMessages,
        unreadNotifications: state.unreadNotifications,
        isLoggedIn: getAuthAlias(rawState, "isLoggedIn")
    };
}

function getState(key?: string | null): any {
    if (key === undefined || key === null || key === "") {
        return buildPublicSnapshot();
    }
    const rootKey = String(key).split(".")[0];
    if (!allowedKeys.has(rootKey)) {
        throw new Error(`Invalid state key:${key}`);
    }
    if (key.includes(".")) {
        return getValueByPath(key);
    }
    return getStateValue(key);
}

function subscribe(key: string, fn: StateListener): () => void {
    const rootKey = String(key).split(".")[0];
    if (!allowedKeys.has(rootKey)) {
        throw new Error(`Cannot subscribe to invalid key:${key}`);
    }
    return stateListenerRegistry.subscribe(key, fn);
}

function unsubscribe(key: string, fn: StateListener): void {
    stateListenerRegistry.unsubscribe(key, fn);
}

function subscribeDeep(path: string, fn: StateListener): () => void {
    return stateListenerRegistry.subscribeDeep(path, fn);
}

function unsubscribeDeep(path: string, fn: StateListener): void {
    stateListenerRegistry.unsubscribeDeep(path, fn);
}

function clearAllListeners(): void {
    stateListenerRegistry.clearAllListeners();
}

function clearState(persist = true): void {
    setState({
        token: null,
        user: null,
        isLoggedIn: false,
        userProfile: {},
        roles: [],
        permissions: [],
        currentChatId: null,
        unreadMessages: 0,
        unreadNotifications: 0
    }, persist);

    rawState.isLoggedIn = false;

    if (persist) {
        removeStorage("token");
        removeStorage("user");
        removeStorage("userProfile");
        removeStorage("roles");
        removeStorage("permissions");
        removeStorage("unreadMessages");
        removeStorage("unreadNotifications");
    }
}

migrateLegacyToken();

export {
    getState,
    setState,
    clearState,
    subscribe,
    unsubscribe,
    clearAllListeners,
    getRouteCache,
    getRouteStateMap,
    readRouteState,
    writeRouteState,
    deleteRouteState,
    clearRouteState,
    saveScroll,
    restoreScroll,
    state
};
