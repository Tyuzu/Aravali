import { API_URL, getState, setState } from "../state/state.js";
import { generateUUID } from "../utils/genUUID.js";

/* =========================================================
   TYPES & INTERFACES
========================================================= */
export interface RawUserRecord {
  id?: string;
  userid?: string;
  username?: string;
  roles?: string | string[];
  role?: string | string[];
  permissions?: string | string[];
  [key: string]: unknown;
}

export interface AuthPayload {
  user: RawUserRecord | null;
  userid: string | null;
  username: string;
  roles: string[];
  permissions: string[];
  auth: {
    isAuthenticated: boolean;
    user: RawUserRecord | null;
    roles: string[];
    permissions: string[];
  };
}

export interface LockResult {
  lockedByOtherTab?: boolean;
}

export interface RefreshLockData {
  owner: string;
  ts: number;
}

/* =========================================================
   CONSTANTS & INITIALIZATION
========================================================= */
const REFRESH_LOCK_TTL = 10_000;
const REFRESH_WAIT_TIMEOUT = 12_000;
const TAB_ID = typeof crypto !== "undefined" && crypto.randomUUID ? crypto.randomUUID() : generateUUID();
const REFRESH_LOCK_KEY = "__refresh_lock__";
const AUTH_CHANNEL = typeof BroadcastChannel !== "undefined" ? new BroadcastChannel("auth_channel") : null;

/* =========================================================
   NAVIGATION REQUEST CANCELLATION
========================================================= */
export let navigationAbortController = new AbortController();

export function abortInflightApiRequests(): void {
  navigationAbortController.abort();
  navigationAbortController = new AbortController();
}

/* =========================================================
   HELPERS
========================================================= */
function normalizeRoles(value: unknown): string[] {
  if (Array.isArray(value)) {
    return [
      ...new Set(
        value
          .filter((role): role is string => role !== null && role !== undefined && String(role).trim() !== "")
          .map((role) => String(role).trim())
      )
    ];
  }
  if (typeof value === "string" && value.trim()) {
    return [value.trim()];
  }
  return [];
}

function normalizePermissions(value: unknown): string[] {
  if (Array.isArray(value)) {
    return [
      ...new Set(
        value
          .filter(
            (permission): permission is string =>
              permission !== null && permission !== undefined && String(permission).trim() !== ""
          )
          .map((permission) => String(permission).trim())
      )
    ];
  }
  if (typeof value === "string" && value.trim()) {
    return [value.trim()];
  }
  return [];
}

/* =========================================================
   REFRESH LOCK
========================================================= */
async function withRefreshLock<T>(taskCallback: () => Promise<T>): Promise<T | LockResult> {
  if (typeof navigator !== "undefined" && navigator.locks) {
    return navigator.locks.request("auth_refresh_lock", async () => {
      return taskCallback();
    });
  }

  const now = Date.now();
  try {
    const raw = localStorage.getItem(REFRESH_LOCK_KEY);
    if (raw) {
      const lock: RefreshLockData = JSON.parse(raw);
      const age = now - (lock.ts || 0);
      if (age < REFRESH_LOCK_TTL && lock.owner !== TAB_ID) {
        return { lockedByOtherTab: true };
      }
    }
    localStorage.setItem(REFRESH_LOCK_KEY, JSON.stringify({ owner: TAB_ID, ts: now }));
  } catch {
    // Continue without cross-tab lock.
  }

  try {
    return await taskCallback();
  } finally {
    try {
      const raw = localStorage.getItem(REFRESH_LOCK_KEY);
      if (raw) {
        const lock: RefreshLockData = JSON.parse(raw);
        if (lock.owner === TAB_ID) {
          localStorage.removeItem(REFRESH_LOCK_KEY);
        }
      }
    } catch {
      // Ignore lock cleanup failures.
    }
  }
}

/* =========================================================
   WAIT FOR ANOTHER TAB
========================================================= */
function waitForSessionUpdate(timeoutMs: number = REFRESH_WAIT_TIMEOUT): Promise<boolean> {
  return new Promise((resolve) => {
    const started = Date.now();
    const timer = setInterval(() => {
      const auth = getState("auth");
      if (auth?.isAuthenticated) {
        clearInterval(timer);
        resolve(true);
        return;
      }
      if (Date.now() - started >= timeoutMs) {
        clearInterval(timer);
        resolve(false);
      }
    }, 100);
  });
}

/* =========================================================
   SESSION / COOKIE REFRESH
========================================================= */
let refreshPromise: Promise<boolean> | null = null;

export async function refreshToken(): Promise<boolean> {
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = (async (): Promise<boolean> => {
    let success = false;

    const lockResult = await withRefreshLock(async () => {
      try {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 10_000);

        const response = await fetch(`${API_URL}/auth/refresh`, {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
            "X-Refresh-Intent": "1"
          },
          signal: controller.signal
        });

        clearTimeout(timeoutId);

        if (!response.ok) {
          success = false;
          return;
        }

        const resData = await response.json().catch(() => null);
        const data = resData?.data && typeof resData.data === "object" ? resData.data : resData;

        const userId =
          resData?.userid ??
          resData?.UserID ??
          data?.userid ??
          data?.UserID ??
          null;

        const username = resData?.username ?? data?.username ?? "";
        const roles = normalizeRoles(
          resData?.roles ?? resData?.role ?? data?.roles ?? data?.role
        );
        const permissions = normalizePermissions(
          resData?.permissions ?? data?.permissions
        );

        const rawUser =
          resData?.user && typeof resData.user === "object"
            ? resData.user
            : data?.user && typeof data.user === "object"
            ? data.user
            : null;

        const authPayload: AuthPayload = {
          user: rawUser,
          userid: userId,
          username,
          roles,
          permissions,
          auth: {
            isAuthenticated: true,
            user: rawUser,
            roles,
            permissions
          }
        };

        setState(authPayload, true);
        AUTH_CHANNEL?.postMessage({
          type: "SESSION_REFRESHED",
          payload: authPayload
        });

        success = true;
      } catch (error: unknown) {
        if (error instanceof Error && error.name === "AbortError") {
          console.warn("[Auth] Session refresh request timed out.");
        } else {
          console.error("[Auth] Session refresh request failed:", error);
        }
        success = false;
      }
    });

    if (lockResult && "lockedByOtherTab" in lockResult && lockResult.lockedByOtherTab) {
      success = await waitForSessionUpdate();
    }

    return success;
  })();

  try {
    return await refreshPromise;
  } finally {
    refreshPromise = null;
  }
}

/* =========================================================
   AUTH CHANNEL
========================================================= */
AUTH_CHANNEL?.addEventListener("message", (event: MessageEvent) => {
  if (event.data?.type === "SESSION_REFRESHED") {
    if (event.data.payload) {
      setState(event.data.payload, true);
    }
  }
  if (event.data?.type === "LOGOUT") {
    window.dispatchEvent(new CustomEvent("auth:remote-logout"));
  }
});

/* =========================================================
   LOCAL AUTH EVENTS
========================================================= */
if (typeof window !== "undefined") {
  window.addEventListener("auth:logout", (event: CustomEventInit) => {
    if (!event.detail?.broadcast) {
      return;
    }
    AUTH_CHANNEL?.postMessage({
      type: "LOGOUT"
    });
  });
}

/* =========================================================
   VISIBILITY REFRESH
========================================================= */
if (typeof document !== "undefined") {
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState !== "visible") {
      return;
    }
    const auth = getState("auth");
    if (auth?.isAuthenticated) {
      refreshToken().then((success) => {
        if (!success) {
          window.dispatchEvent(new CustomEvent("auth:unauthorized"));
        }
      });
    }
  });
}