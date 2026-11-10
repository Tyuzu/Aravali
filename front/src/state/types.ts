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
