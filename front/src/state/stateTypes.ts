import { apiConfig } from "../config/env.js";
import { saveScroll, restoreScroll } from "./storage.js";

/* =========================================================
    TYPES & INTERFACES
========================================================= */
export interface User {
    id?: string | number;
    userid?: string | number;
    username?: string;
    name?: string;
    [key: string]: any;
}

export interface AuthState {
    isAuthenticated: boolean;
    loading: boolean;
    accessToken: string | null;
    user: User | null;
    roles: string[];
    permissions: string[];
}

export interface AppState {
    auth: AuthState;
    userProfile: Record<string, any>;
    socket: any | null;
    environment: Record<string, any>;
    lang: string;
    lastPath: string;
    currentRoute: any | null;
    routeCache: Map<any, any>;
    routeState: Map<any, any>;
    currentChatId: string | number | null;
    isLoading: boolean;
    unreadMessages: number;
    unreadNotifications: number;
    isLoggedIn: boolean;
    [key: string]: any;
}

export type StateListener = (value: any, state: AppState) => void;

/* =========================================================
    API CONFIG EXPORTS
========================================================= */
export const {
    MAIN_URL, EMBED_URL, BANNERDROP_URL, API_URL, STRIPE_URL, AD_URL, SEARCH_URL, MERE_URL, MERE_WS, CHAT_URL, CHAT_WS, MUSIC_URL, LIVE_URL, SRC_URL, FILEDROP_URL, CHATDROP_URL 
} = apiConfig;

// Re-export scroll utilities
export { saveScroll, restoreScroll };