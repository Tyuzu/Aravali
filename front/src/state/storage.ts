import { apiConfig } from "../config/env.js";

/* =========================================================
    STORAGE CONFIG & KEYS
========================================================= */
export const PERSISTED_KEYS = new Set<string>([
    "userProfile", "user", "roles", "permissions",
    "favFarms",
    "unreadMessages", "unreadNotifications"
]);

export const SESSION_KEYS = new Set<string>(["token"]);

/* =========================================================
    STORAGE READERS & WRITERS
========================================================= */
export function readSessionStorage(key: string): string | null {
    try {
        return sessionStorage.getItem(key);
    } catch {
        return null;
    }
}

export function readLocalStorage(key: string): string | null {
    try {
        return localStorage.getItem(key);
    } catch {
        return null;
    }
}

export function readStorage(key: string): string | null {
    if (SESSION_KEYS.has(key)) {
        return readSessionStorage(key);
    }
    return readLocalStorage(key);
}

export function serializeValue(value: any): string | null {
    if (typeof value === "string") {
        return value;
    }
    try {
        return JSON.stringify(value);
    } catch (error) {
        console.warn(`[STATE] Unable to serialize state key "${String(error)}":`, error);
        return null;
    }
}

export function writeSessionStorage(key: string, value: any): boolean {
    try {
        if (value === null || value === undefined) {
            sessionStorage.removeItem(key);
            return true;
        }
        const serialized = serializeValue(value);
        if (serialized === null) {
            return false;
        }
        sessionStorage.setItem(key, serialized);
        return true;
    } catch (error) {
        console.warn(`[STATE] Failed writing session key "${key}":`, error);
        return false;
    }
}

export function writeLocalStorage(key: string, value: any): boolean {
    try {
        if (value === null || value === undefined) {
            localStorage.removeItem(key);
            return true;
        }
        const serialized = serializeValue(value);
        if (serialized === null) {
            return false;
        }
        localStorage.setItem(key, serialized);
        return true;
    } catch (error) {
        console.warn(`[STATE] Failed writing persistent key "${key}":`, error);
        return false;
    }
}

export function writeStorage(key: string, value: any): boolean {
    if (SESSION_KEYS.has(key)) {
        return writeSessionStorage(key, value);
    }
    if (PERSISTED_KEYS.has(key)) {
        return writeLocalStorage(key, value);
    }
    return false;
}

export function removeStorage(key: string): void {
    try {
        sessionStorage.removeItem(key);
    } catch {
        // Ignore
    }
    try {
        localStorage.removeItem(key);
    } catch {
        // Ignore
    }
}

export function safeParseFromStorage<T = any>(key: string, fallback: T | null = null): T | string | null {
    const raw = readStorage(key);
    if (raw === null || raw === "") {
        return fallback;
    }
    try {
        return JSON.parse(raw);
    } catch {
        return raw;
    }
}

export function readPersistentJSON<T = any>(key: string, fallback: T | null = null): T | null {
    const raw = readLocalStorage(key);
    if (raw === null || raw === "") {
        return fallback;
    }
    try {
        return JSON.parse(raw);
    } catch {
        return fallback;
    }
}

export function readPersistentNumber(key: string, fallback = 0): number {
    const value = readPersistentJSON(key, fallback);
    const number = Number(value);
    return Number.isFinite(number) ? number : fallback;
}

/* =========================================================
    LEGACY TOKEN MIGRATION
========================================================= */
export function migrateLegacyToken(): void {
    const sessionToken = readSessionStorage("token");
    const localToken = readLocalStorage("token");
    if (!sessionToken && localToken) {
        try {
            sessionStorage.setItem("token", localToken);
        } catch (error) {
            console.warn("[AUTH] Unable to migrate legacy token:", error);
        }
    }
    try {
        localStorage.removeItem("token");
    } catch {
        // Ignore.
    }
}

/* =========================================================
    ROUTE CACHE & SCROLL STATE
========================================================= */
export const routeCache = new Map<any, any>();
export const routeState = new Map<any, any>();
const scrollPositions = new Map<any, { top: number; left: number }>();

export function saveScroll(container: HTMLElement | null, location: any): void {
    if (!container) return;
    scrollPositions.set(location, {
        top: container.scrollTop || window.scrollY || 0,
        left: container.scrollLeft || window.scrollX || 0
    });
}

export function restoreScroll(container: HTMLElement | null, location: any): void {
    if (!container) return;
    const pos = scrollPositions.get(location);
    if (pos) {
        if (container === document.body || container === document.documentElement) {
            window.scrollTo(pos.left, pos.top);
        } else {
            container.scrollTop = pos.top;
            container.scrollLeft = pos.left;
        }
    } else {
        container.scrollTop = 0;
        window.scrollTo(0, 0);
    }
}