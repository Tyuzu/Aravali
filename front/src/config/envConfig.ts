/**
 * Environment Configuration
 * Vite + Netlify same-origin API configuration.
 */

export const webSiteName: string = "Aravali";

// Access Vite environment variables safely without TypeScript compiler errors
const env = (import.meta as unknown as { env: Record<string, string | undefined> }).env || {};

// Normalize base URLs by removing trailing slashes
const MAIN_URL: string = (env.VITE_MAIN_URL || "").replace(/\/+$/, "");
const BANNERDROP_URL: string = (env.VITE_BANNERDROP_URL || "").replace(/\/+$/, "");
const MODE: string = env.MODE || "development";

/**
 * Safely construct relative or absolute URLs.
 */
const buildURL = (baseURL: string, path: string): string => {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  return baseURL ? `${baseURL}${normalizedPath}` : normalizedPath;
};

/**
 * Safely compute WebSocket URL with SSR & Relative Path fallbacks.
 */
const getWebSocketURL = (): string => {
  // 1. Explicit full API URL configured (e.g., http://localhost:4000)
  if (MAIN_URL) {
    try {
      const url = new URL(MAIN_URL, typeof window !== "undefined" ? window.location.origin : "http://localhost");
      url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
      return url.toString().replace(/\/+$/, "");
    } catch {
      // Fallback for invalid URLs during build
      return MAIN_URL.replace(/^http/, "ws");
    }
  }

  // 2. Same-origin browser runtime
  if (typeof window !== "undefined") {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    return `${protocol}//${window.location.host}`;
  }

  // 3. Fallback for Node/SSR environments
  return "";
};

export function normalizeSocketBase(baseUrl: string): string {
  const cleanBase = (baseUrl || "").replace(/\/+$/, "");

  if (!cleanBase) {
    return "";
  }

  if (/^wss?:\/\//i.test(cleanBase)) {
    return cleanBase;
  }

  if (/^https?:\/\//i.test(cleanBase)) {
    return cleanBase.replace(/^http/, "ws");
  }

  if (typeof window !== "undefined") {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    return `${protocol}//${cleanBase.replace(/^\/\//, "")}`;
  }

  return cleanBase.replace(/^http/, "ws");
}

export function buildWebSocketUrl(
  baseUrl: string,
  path: string = "",
  token?: string | null
): string {
  const normalizedBase = normalizeSocketBase(baseUrl || "");
  const safePath = path ? (path.startsWith("/") ? path : `/${path}`) : "";
  const url = `${normalizedBase}${safePath}`;

  if (!token) {
    return url;
  }

  const separator = url.includes("?") ? "&" : "?";
  return `${url}${separator}token=${encodeURIComponent(token)}`;
}

const WS_URL: string = getWebSocketURL();

export interface ApiConfig {
  MAIN_URL: string;
  BANNERDROP_URL: string;
  API_URL: string;
  STRIPE_URL: string;
  AD_URL: string;
  SEARCH_URL: string;
  MERE_URL: string;
  MUSIC_URL: string;
  LIVE_URL: string;
  EMBED_URL: string;
  MERE_WS: string;
  CHAT_URL: string;
  CHAT_WS: string;
  SRC_URL: string;
  FILEDROP_URL: string;
  CHATDROP_URL: string;
  isDev: boolean;
  isStaging: boolean;
  isProduction: boolean;
  environment: string;
}

/**
 * Centralized API configuration.
 */
export const apiConfig: ApiConfig = {
  /* Base URLs */
  MAIN_URL,
  BANNERDROP_URL,

  /* API Endpoints */
  API_URL: buildURL(MAIN_URL, "/api/v1"),
  STRIPE_URL: buildURL(MAIN_URL, "/api/v1/stripe"),
  AD_URL: buildURL(MAIN_URL, "/api/sda"),
  SEARCH_URL: buildURL(MAIN_URL, "/api/v1"),
  MERE_URL: buildURL(MAIN_URL, "/api/v1"),
  MUSIC_URL: buildURL(MAIN_URL, "/api/v1"),
  LIVE_URL: buildURL(MAIN_URL, "/api/v1"),
  EMBED_URL: buildURL(MAIN_URL, "/embed"),

  /* WebSockets */
  MERE_WS: WS_URL,
  CHAT_URL: MAIN_URL,
  CHAT_WS: WS_URL ? buildWebSocketUrl(WS_URL, "/ws/newchat/chat") : "",

  /* Media & Static */
  SRC_URL: buildURL(BANNERDROP_URL, "/static"),
  FILEDROP_URL: buildURL(BANNERDROP_URL, "/api/v1/filedrop"),
  CHATDROP_URL: buildURL(BANNERDROP_URL, "/api/v1/filedrop"),

  /* Mode Flags */
  isDev: MODE === "development" || MODE === "dev",
  isStaging: MODE === "staging",
  isProduction: MODE === "production",
  environment: MODE,
};

export default apiConfig;