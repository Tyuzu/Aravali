import "../../css/ui/createTabs.css";
import { createElement, ElementAttributes } from "../components/createElement.js";
import { makeDraggableScroll } from "../components/dragnav.js";
import { readRouteState, writeRouteState } from "../state/state.js";

/* =========================================================
   TYPES & INTERFACES
========================================================= */

export interface TabItem {
  id: string;
  title: string;
  render: (container: HTMLElement) => void;
}

interface RouteTabState {
  activeTab?: string;
  [key: string]: unknown;
}

export type OnTabChangeCallback = (tabId: string) => void;

/* =========================================================
   ROUTE STATE HELPERS
========================================================= */

function getRouteState(locationKey: string): RouteTabState | null {
  const saved = readRouteState(locationKey);
  return saved && typeof saved === "object" ? (saved as RouteTabState) : null;
}

function setRouteState(locationKey: string, value: RouteTabState): void {
  writeRouteState(locationKey, value);
}

/* =========================================================
   EXPORTED TAB FUNCTIONS
========================================================= */

/**
 * High-order layout helper that handles persisting active tab views via localStorage.
 */
export function persistTabs(
  container: HTMLElement | null,
  tabs: TabItem[],
  storageKey: string | null = null
): void {
  if (!container) return;
  try {
    const resolvedKey = storageKey || "default-tabs";
    const activeTabId = localStorage.getItem(resolvedKey);

    const tabsElement = createTabs(
      tabs,
      resolvedKey,
      activeTabId,
      (newTabId) => {
        localStorage.setItem(resolvedKey, newTabId);
        const tabState = getRouteState(resolvedKey) || {};
        tabState.activeTab = newTabId;
        setRouteState(resolvedKey, tabState);
      }
    );

    container.appendChild(tabsElement);
  } catch (err) {
    console.warn("Tabs component failed to initialize in container:", err);
  }
}

/**
 * Generates an accessible, performance-lazy rendering tab node interface with WAI-ARIA keyboard navigation.
 */
export function createTabs(
  tabs: TabItem[],
  routeKey: string | null = null,
  initialTabId: string | null = null,
  onTabChange: OnTabChangeCallback | null = null
): HTMLDivElement {
  if (!Array.isArray(tabs) || tabs.length === 0) {
    return createElement("div", { class: "tabs-empty" });
  }

  const tabContainer = createElement("div", { class: "tabs-container" });
  const tabButtons = createElement("div", {
    class: "tab-buttons",
    role: "tablist",
    "aria-label": "Content Navigation Tabs",
  });
  const tabContents = createElement("div", { class: "tab-contents" });

  const tabContentMap = new Map<string, HTMLElement>();
  const buttonMap = new Map<string, HTMLElement>();
  const tabIds = tabs.map((t) => t.id);

  // State Resolution: Route state takes precedence, followed by initialTabId, then fallback to first tab
  let initial = tabs[0]?.id;
  if (routeKey) {
    const saved = getRouteState(routeKey);
    if (saved?.activeTab && tabIds.includes(saved.activeTab)) {
      initial = saved.activeTab;
    } else if (initialTabId && tabIds.includes(initialTabId)) {
      initial = initialTabId;
    }
  } else if (initialTabId && tabIds.includes(initialTabId)) {
    initial = initialTabId;
  }

  // Instantiate Interactive Tab View Models
  tabs.forEach(({ id, title }) => {
    const panelId = `panel-${id}`;
    const buttonId = `tab-btn-${id}`;
    const isActive = id === initial;

    const contentContainer = createElement("article", {
      id: panelId,
      class: `tab-content${isActive ? " active" : ""}`,
      role: "tabpanel",
      "aria-labelledby": buttonId,
      tabindex: 0,
    });

    const buttonAttributes: ElementAttributes = {
      id: buttonId,
      class: `tab-button${isActive ? " active" : ""}`,
      role: "tab",
      "aria-controls": panelId,
      "aria-selected": isActive ? "true" : "false",
      tabindex: isActive ? 0 : -1,
      dataset: { id },
      events: {
        click: () => activateTab(id, true),
        keydown: ((e: Event) => {
          const keyEvent = e as KeyboardEvent;
          const currentIndex = tabIds.indexOf(id);
          let targetId: string | null = null;

          if (keyEvent.key === "Enter" || keyEvent.key === " ") {
            keyEvent.preventDefault();
            activateTab(id, true);
            return;
          }

          if (keyEvent.key === "ArrowRight" || keyEvent.key === "ArrowDown") {
            keyEvent.preventDefault();
            const nextIndex = (currentIndex + 1) % tabIds.length;
            targetId = tabIds[nextIndex] ?? null;
          } else if (keyEvent.key === "ArrowLeft" || keyEvent.key === "ArrowUp") {
            keyEvent.preventDefault();
            const prevIndex = (currentIndex - 1 + tabIds.length) % tabIds.length;
            targetId = tabIds[prevIndex] ?? null;
          } else if (keyEvent.key === "Home") {
            keyEvent.preventDefault();
            targetId = tabIds[0] ?? null;
          } else if (keyEvent.key === "End") {
            keyEvent.preventDefault();
            targetId = tabIds[tabIds.length - 1] ?? null;
          }

          if (targetId) {
            activateTab(targetId, true);
          }
        }) as EventListener,
      },
    };

    const tabButton = createElement("div", buttonAttributes, title);

    tabButtons.appendChild(tabButton);
    tabContents.appendChild(contentContainer);

    tabContentMap.set(id, contentContainer);
    buttonMap.set(id, tabButton);
  });

  // Attach horizontal swipe nav features if active
  if (typeof makeDraggableScroll === "function") {
    makeDraggableScroll(tabButtons);
  }

  tabContainer.appendChild(tabButtons);
  tabContainer.appendChild(tabContents);

  /**
   * Toggles layout active states and handles lazy evaluation injections
   */
  function activateTab(tabId: string, focusManagement = false): void {
    tabs.forEach(({ id, render }) => {
      const btn = buttonMap.get(id);
      const content = tabContentMap.get(id);
      const isActive = id === tabId;

      if (btn && content) {
        btn.classList.toggle("active", isActive);
        btn.setAttribute("aria-selected", isActive ? "true" : "false");
        btn.setAttribute("tabindex", isActive ? "0" : "-1");
        content.classList.toggle("active", isActive);

        // Lazy Rendering Verification Block
        if (isActive && !content.dataset.rendered) {
          if (typeof render === "function") {
            render(content);
          }
          content.dataset.rendered = "true";
        }
      }
    });

    if (focusManagement) {
      const activeBtn = buttonMap.get(tabId);
      activeBtn?.focus();
    }

    // Router history persistence integration check
    if (routeKey) {
      const tabState = getRouteState(routeKey) || {};
      tabState.activeTab = tabId;
      setRouteState(routeKey, tabState);
    }

    if (typeof onTabChange === "function") {
      onTabChange(tabId);
    }
  }

  // Trigger initial lazy render on startup
  if (initial) {
    const initialContent = tabContentMap.get(initial);
    if (initialContent && !initialContent.dataset.rendered) {
      const tabObj = tabs.find((t) => t.id === initial);
      if (tabObj && typeof tabObj.render === "function") {
        queueMicrotask(() => {
          tabObj.render(initialContent);
          initialContent.dataset.rendered = "true";
        });
      }
    }
  }

  return tabContainer;
}