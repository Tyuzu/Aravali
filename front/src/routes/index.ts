import { createheader } from "../components/layout/header.js";
import { createNav, highlightActiveNav } from "../components/layout/navigation.js";
import { Footer } from "../components/layout/footer.js";
import LoadingSpinner from "../components/ui/LoadingSpinner.js";
import { render } from "./router.js";
import { getState, setState, restoreScroll, subscribe } from "../state/state.js";
import { track } from "../services/activity/metrics.js";
import { startPerfMonitoring } from "../services/activity/perfMonitor.js";
import { normalizeStoredRedirect } from "../utils/navPath.js";
import {
    getCurrentAppLocation,
    parseAppLocation,
    toBrowserTarget,
    navigate,
    registerContentLoader
} from "./navigate.js";

/* Re-export navigation utilities for external modules */
export { getCurrentAppLocation, parseAppLocation, navigate } from "./navigate.js";

/* =========================================================
   TYPES & INTERFACES
========================================================= */
interface LayoutElements {
    header: HTMLElement | null;
    nav: HTMLElement | null;
    main: HTMLElement | null;
    footer: HTMLElement | null;
}

interface RouteRedirectResult {
    redirect?: string;
    [key: string]: unknown;
}

/* =========================================================
   ROUTE STATE HELPER
========================================================= */
function getRouteStateHelper(location: string): unknown {
    const state = getState() || {};
    const routeStates = state.routeState || state.scrollPositions || {};
    return routeStates[location] || null;
}

/* =========================================================
   LAYOUT STATE & ELEMENTS CACHE
========================================================= */
const layoutState = {
    isHydrated: false,
    headerRendered: false,
    navRendered: false,
    footerRendered: false
};

const elements: LayoutElements = {
    header: null,
    nav: null,
    footer: null,
    main: null
};

let pageLoadingIndicator: HTMLElement | null = null;

function showPageLoadingIndicator(): void {
    if (pageLoadingIndicator && pageLoadingIndicator.isConnected) {
        return;
    }

    const indicator = document.createElement("div");
    indicator.id = "page-navigation-loader";
    indicator.className = "page-navigation-loader";
    indicator.setAttribute("role", "status");
    indicator.setAttribute("aria-live", "polite");
    indicator.setAttribute("aria-busy", "true");
    indicator.appendChild(LoadingSpinner.createElement());

    const root = document.body || document.documentElement;
    if (root) {
        root.appendChild(indicator);
        pageLoadingIndicator = indicator;
    }
}

function hidePageLoadingIndicator(): void {
    if (!pageLoadingIndicator) return;

    pageLoadingIndicator.remove();
    pageLoadingIndicator = null;
}

function getElements(): LayoutElements {
    if (!elements.main) {
        elements.header = document.getElementById("pageheader");
        elements.nav = document.getElementById("primary-nav");
        elements.main = document.getElementById("content");
        elements.footer = document.getElementById("pagefooter");
    }
    return elements;
}

/* =========================================================
   LAYOUT RENDERERS (UNIFIED)
========================================================= */
function renderHeader(force = false): void {
    const { header } = getElements();
    if (!header) return;

    if (!layoutState.headerRendered || force) {
        createheader();
        layoutState.headerRendered = true;
    }
}

function renderNav(force = false): void {
    const { nav } = getElements();
    if (!nav) return;

    if (!layoutState.navRendered || force) {
        const nextNav = createNav();
        if (nextNav) nav.replaceChildren(nextNav);
        layoutState.navRendered = true;
    }
}

function renderFooter(force = false): void {
    const { footer } = getElements();
    if (!footer) return;

    if (!layoutState.footerRendered || force) {
        const nextFooter = Footer();
        if (nextFooter) footer.replaceChildren(nextFooter);
        layoutState.footerRendered = true;
    }
}

function updateLayout(force = false): void {
    renderHeader(force);
    renderNav(force);
    renderFooter(force);
}

/* =========================================================
   AUTH HYDRATION
========================================================= */
function getStoredValue(key: string): string | null {
    return sessionStorage.getItem(key) || localStorage.getItem(key) || null;
}

export function hydrateAuthState(force = false): void {
    if (layoutState.isHydrated && !force) return;

    const state = getState() || {};
    const auth = state.auth || {};
    const token = getStoredValue("token");
    const userRaw = getStoredValue("user");
    let user: unknown = null;

    if (userRaw) {
        try {
            const trimmed = userRaw.trim();
            user = (trimmed.startsWith("{") || trimmed.startsWith("[")) 
                ? JSON.parse(trimmed) 
                : userRaw;
        } catch (error) {
            console.warn("Failed parsing stored user JSON:", error);
            user = userRaw;
        }
    }

    setState(
        {
            token,
            user,
            auth: {
                ...auth,
                accessToken: token || auth.accessToken || null,
                isAuthenticated: Boolean(token || auth.accessToken),
                user: user || auth.user || null
            }
        },
        true
    );

    layoutState.isHydrated = true;
}

/* =========================================================
   CONTENT LOADING & ROUTING
========================================================= */
function renderSkeleton(): string {
    return `
        <div class="skeleton-loader" aria-live="polite" aria-busy="true">
            <div class="skeleton-bar short"></div>
            <div class="skeleton-bar"></div>
            <div class="skeleton-bar long"></div>
        </div>
    `;
}

export async function loadContent(
    rawLocation: string = getCurrentAppLocation(),
    { redirectDepth = 0 }: { redirectDepth?: number } = {}
): Promise<unknown> {
    const { main } = getElements();
    if (!main) {
        throw new Error("SPA content container #content was not found.");
    }

    showPageLoadingIndicator();

    try {
        hydrateAuthState();
        updateLayout(false);

        main.classList.add("is-transitioning");
        main.innerHTML = renderSkeleton();

        await new Promise((resolve) => {
            requestAnimationFrame(() => {
                window.setTimeout(resolve, 80);
            });
        });

        const parsed = parseAppLocation(rawLocation);
        const result = await render(parsed.fullPath, main);

        /* --- ROUTER REDIRECT HANDLING --- */
        const routeResult = result as RouteRedirectResult;
        if (routeResult && typeof routeResult.redirect === "string" && routeResult.redirect) {
            if (redirectDepth >= 10) {
                throw new Error("Too many consecutive route redirects.");
            }

            const browserTarget = toBrowserTarget(routeResult.redirect, parsed.mode);
            history.replaceState(null, "", browserTarget);

            track("route_redirect", {
                from: parsed.fullPath,
                to: routeResult.redirect
            });

            return loadContent(getCurrentAppLocation(), {
                redirectDepth: redirectDepth + 1
            });
        }

        /* --- NAVIGATION UI --- */
        highlightActiveNav(parsed.path);

        /* --- SCROLL RESTORATION --- */
        const routeState = getRouteStateHelper(getCurrentAppLocation());
        if (routeState) {
            requestAnimationFrame(() => {
                restoreScroll(main, routeState);
            });
        }

        requestAnimationFrame(() => {
            main.classList.remove("is-transitioning");
        });

        return result;
    } finally {
        hidePageLoadingIndicator();
        requestAnimationFrame(() => {
            main.classList.remove("is-transitioning");
        });
    }
}

// Register loader with navigate module to avoid circular dependency
registerContentLoader((location) => loadContent(location));

/* =========================================================
   INITIAL RENDER API
========================================================= */
export async function renderPage(): Promise<void> {
    startPerfMonitoring();
    await loadContent(getCurrentAppLocation());
}


/* =========================================================
   AUTH REACTIVITY & CROSS-TAB SYNC
========================================================= */
subscribe("token", (token: unknown) => {
    updateLayout(true);

    if (!token) return;

    const current = parseAppLocation(getCurrentAppLocation());
    if (current.path !== "/login") return;

    const redirect = normalizeStoredRedirect(sessionStorage.getItem("redirectAfterLogin"));
    sessionStorage.removeItem("redirectAfterLogin");

    const target = redirect || "/";

    queueMicrotask(() => {
        navigate(target, { replace: true }).catch((error) => {
            console.error("Post-login redirect failed:", error);
        });
    });
});

subscribe("lang", async () => {
    try {
        updateLayout(true);
        await loadContent(getCurrentAppLocation());
    } catch (error) {
        console.error("[LANG] Failed to refresh UI after language change:", error);
    }
});

window.addEventListener("storage", (event: StorageEvent) => {
    if (event.key === "token" || event.key === "user") {
        hydrateAuthState(true);
        updateLayout(true);
    }

    if (event.key === "lang") {
        updateLayout(true);
        loadContent(getCurrentAppLocation()).catch((error) => {
            console.error("[LANG] Failed to re-render after storage language change:", error);
        });
    }
});