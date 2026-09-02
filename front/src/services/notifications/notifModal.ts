import "../../../css/subpages/notifications.css";
import Modal from "../../components/ui/Modal.js";
import { createElement } from "../../components/createElement.js";
import {
  getNotifications,
  markNotificationAsRead,
  markAllNotificationsAsRead,
  clearAllNotifications,
  NotificationItem,
} from "./notifService.js";
import * as idxDB from "../../utils/idxDB.js";

interface SystemLog {
  id?: string | number;
  notificationid?: string | number;
  type?: "info" | "error" | "success" | string;
  title?: string;
  message?: string;
  createdAt?: string | number | Date;
  isRead?: boolean;
}

type NotificationFilter = "all" | "unread";
type SystemFilter = "all" | "unread" | "error" | "success";

const UI_STATE = {
  search: "",
  activityFilter: "all" as NotificationFilter,
  systemFilter: "all" as SystemFilter,
};

/**
 * Formats a given date string/timestamp into relative human-readable time.
 */
function timeAgo(dateInput: string | number | Date): string {
  const date = new Date(dateInput);
  if (isNaN(date.getTime())) return "";

  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
  const formatter = new Intl.RelativeTimeFormat("en", { numeric: "auto" });

  if (seconds < 60) return "Just now";
  if (seconds < 3600) return formatter.format(-Math.floor(seconds / 60), "minute");
  if (seconds < 86400) return formatter.format(-Math.floor(seconds / 3600), "hour");
  if (seconds < 2592000) return formatter.format(-Math.floor(seconds / 86400), "day");

  return date.toLocaleDateString();
}

function matchesQuery(value: string | undefined, query: string): boolean {
  if (!query) return true;
  return (value || "").toLowerCase().includes(query.toLowerCase());
}

function filterNotifications(items: NotificationItem[], filter: NotificationFilter, query: string): NotificationItem[] {
  return items.filter((item) => {
    const matchesFilter = filter === "all" ? true : !item.isRead;
    const text = [item.title, item.type, item.message].join(" ");
    return matchesFilter && matchesQuery(text, query);
  });
}

function filterSystemLogs(items: SystemLog[], filter: SystemFilter, query: string): SystemLog[] {
  return items.filter((item) => {
    const type = String(item.type || "info").toLowerCase();
    const matchesFilter =
      filter === "all"
        ? true
        : filter === "unread"
          ? !item.isRead
          : type === filter;
    const text = [item.title, item.message, type].join(" ");
    return matchesFilter && matchesQuery(text, query);
  });
}

function getUserId(): string | null {
  try {
    const userStr = localStorage.getItem("user");
    if (!userStr) return null;
    const user = JSON.parse(userStr);
    return user.id || user._id || userStr;
  } catch {
    return localStorage.getItem("user");
  }
}

function renderSummaryChips(container: HTMLElement, items: Array<{ label: string; value: number; tone?: "default" | "success" | "error" | "warning" }>): void {
  const row = createElement("div", { class: "notification-summary" });

  items.forEach(({ label, value, tone = "default" }) => {
    row.appendChild(
      createElement("span", { class: `notification-summary__chip notification-summary__chip--${tone}` }, [
        `${label}: ${value}`,
      ])
    );
  });

  container.innerHTML = "";
  container.appendChild(row);
}

export async function openNotificationsModal(): Promise<void> {
  const userId = getUserId();
  let activeTab: "activity" | "system" = "activity";

  const content = createElement("div", { class: "notification-modal" });

  const tabHeader = createElement("div", { class: "notification-tab-header" });
  const activityTabBtn = createElement("button", { class: "notification-tab-button is-active" }, ["Activity"]);
  const systemTabBtn = createElement("button", { class: "notification-tab-button" }, ["System Logs"]);

  const toolbar = createElement("div", { class: "notification-toolbar" });
  const searchInput = createElement("input", {
    class: "notification-search",
    type: "search",
    placeholder: "Search activity...",
    value: UI_STATE.search,
  }) as HTMLInputElement;

  const filterGroup = createElement("div", { class: "notification-filter-group" });
  const filterButtons = {
    all: createElement("button", { class: "notification-filter-btn is-active", type: "button" }, ["All"]),
    unread: createElement("button", { class: "notification-filter-btn", type: "button" }, ["Unread"]),
  };

  Object.entries(filterButtons).forEach(([key, button]) => {
    button.addEventListener("click", () => {
      const nextFilter = key as NotificationFilter;
      UI_STATE.activityFilter = nextFilter;
      applyActivityFilterState();
      renderActivityTab();
    });
    filterGroup.appendChild(button);
  });

  toolbar.appendChild(searchInput);
  toolbar.appendChild(filterGroup);

  tabHeader.appendChild(activityTabBtn);
  tabHeader.appendChild(systemTabBtn);
  content.appendChild(tabHeader);
  content.appendChild(toolbar);

  const tabContentView = createElement("div", { class: "notification-tab-content" });
  content.appendChild(tabContentView);

  Modal({
    title: "📬 Notifications & Logs",
    content: content,
    size: "medium",
    showCloseButton: true,
  });

  searchInput.addEventListener("input", (event: Event) => {
    const value = (event.target as HTMLInputElement).value;
    UI_STATE.search = value;
    if (activeTab === "activity") {
      renderActivityTab();
    } else {
      renderSystemTab();
    }
  });

  activityTabBtn.addEventListener("click", () => {
    if (activeTab === "activity") return;
    activeTab = "activity";
    updateTabStyles();
    applyActivityFilterState();
    renderActivityTab();
  });

  systemTabBtn.addEventListener("click", () => {
    if (activeTab === "system") return;
    activeTab = "system";
    updateTabStyles();
    applySystemFilterState();
    renderSystemTab();
  });

  function applyActivityFilterState(): void {
    const isAll = UI_STATE.activityFilter === "all";
    filterButtons.all.classList.toggle("is-active", isAll);
    filterButtons.unread.classList.toggle("is-active", !isAll);
    searchInput.placeholder = "Search activity...";
  }

  function applySystemFilterState(): void {
    const systemFilters = {
      all: createElement("button", { class: "notification-filter-btn is-active", type: "button" }, ["All"]),
      unread: createElement("button", { class: "notification-filter-btn", type: "button" }, ["Unread"]),
      error: createElement("button", { class: "notification-filter-btn", type: "button" }, ["Errors"]),
      success: createElement("button", { class: "notification-filter-btn", type: "button" }, ["Success"]),
    };

    filterGroup.innerHTML = "";

    Object.entries(systemFilters).forEach(([key, button]) => {
      button.addEventListener("click", () => {
        UI_STATE.systemFilter = key as SystemFilter;
        applySystemFilterState();
        renderSystemTab();
      });
      filterGroup.appendChild(button);
    });

    const active = UI_STATE.systemFilter;
    systemFilters.all.classList.toggle("is-active", active === "all");
    systemFilters.unread.classList.toggle("is-active", active === "unread");
    systemFilters.error.classList.toggle("is-active", active === "error");
    systemFilters.success.classList.toggle("is-active", active === "success");
    searchInput.placeholder = "Search system logs...";
  }

  function updateTabStyles(): void {
    const isActivity = activeTab === "activity";
    activityTabBtn.classList.toggle("is-active", isActivity);
    systemTabBtn.classList.toggle("is-active", !isActivity);
  }

  async function renderActivityTab(): Promise<void> {
    const summaryHost = createElement("div");
    tabContentView.innerHTML = "";
    tabContentView.appendChild(summaryHost);
    summaryHost.innerHTML = `<div class="notification-loading">Loading activity...</div>`;

    try {
      const notifications = (await getNotifications()) || [];
      const filtered = filterNotifications(notifications, UI_STATE.activityFilter, UI_STATE.search);

      const unreadCount = notifications.filter((n) => !n.isRead).length;
      renderSummaryChips(summaryHost, [
        { label: "Total", value: notifications.length },
        { label: "Unread", value: unreadCount, tone: "warning" },
      ]);

      if (!filtered.length) {
        renderEmptyState(tabContentView, "No matching activity updates yet.");
        return;
      }

      const listContainer = createElement("div", { class: "notification-list" });
      filtered.forEach((notification) => {
        listContainer.appendChild(createNotificationCard(notification, userId, renderActivityTab));
      });
      tabContentView.appendChild(listContainer);

      const actionBar = createActionBar(userId, filtered, renderActivityTab);
      if (actionBar) tabContentView.appendChild(actionBar);
    } catch (err) {
      console.error("Failed to load activity notifications:", err);
      tabContentView.innerHTML = `<div class="notification-error">Failed to load activity notifications.</div>`;
    }
  }

  async function renderSystemTab(): Promise<void> {
    const summaryHost = createElement("div");
    tabContentView.innerHTML = "";
    tabContentView.appendChild(summaryHost);
    summaryHost.innerHTML = `<div class="notification-loading">Loading system logs...</div>`;

    let logs: SystemLog[] = [];
    try {
      logs = (await idxDB.getAll()) || [];
      logs.sort((a, b) => new Date(b.createdAt || 0).getTime() - new Date(a.createdAt || 0).getTime());
    } catch (err) {
      console.error("Failed to fetch system logs from IndexedDB:", err);
    }

    const filtered = filterSystemLogs(logs, UI_STATE.systemFilter, UI_STATE.search);
    const unreadCount = logs.filter((log) => !log.isRead).length;
    const errorCount = logs.filter((log) => String(log.type || "info").toLowerCase() === "error").length;
    renderSummaryChips(summaryHost, [
      { label: "Total", value: logs.length },
      { label: "Unread", value: unreadCount, tone: "warning" },
      { label: "Errors", value: errorCount, tone: "error" },
    ]);

    tabContentView.innerHTML = "";
    tabContentView.appendChild(summaryHost);

    if (!filtered.length) {
      renderEmptyState(tabContentView, "No matching system logs found.");
      return;
    }

    const listContainer = createElement("div", { class: "notification-list" });
    filtered.forEach((log) => {
      listContainer.appendChild(createSystemLogCard(log, renderSystemTab));
    });
    tabContentView.appendChild(listContainer);

    const actionBar = createSystemActionBar(filtered, renderSystemTab);
    if (actionBar) tabContentView.appendChild(actionBar);
  }

  applyActivityFilterState();
  updateTabStyles();
  renderActivityTab();
}

function renderEmptyState(container: HTMLElement, message: string): void {
  container.appendChild(
    createElement("div", { class: "notification-empty-state" }, [
      createElement("p", { class: "notification-empty-state__title" }, ["🔔 All caught up"]),
      createElement("p", { class: "notification-empty-state__message" }, [message]),
    ])
  );
}

function createNotificationCard(n: NotificationItem, userId: string | null, onChange?: () => void): HTMLElement {
  let isRead = Boolean(n.isRead);

  const leftContent = createElement("div", { class: "notification-card__content" }, [
    createElement("strong", { class: `notification-card__title${isRead ? " is-read" : ""}` }, [n.title || n.type || "Notification"]),
    createElement("p", { class: "notification-card__message" }, [n.message || "No details provided."]),
    createElement("small", { class: "notification-card__time" }, [timeAgo(n.createdAt || Date.now())]),
  ]);

  const markReadBtn = createElement(
    "button",
    {
      class: "notification-action-btn",
      events: {
        click: async (e: Event) => {
          const mouseEvent = e as MouseEvent;
          mouseEvent.stopPropagation();
          const notifId = n.notificationid || n.id;

          isRead = true;
          card.classList.add("is-read");
          markReadBtn.style.display = "none";

          try {
            await markNotificationAsRead(notifId);
            if (onChange) onChange();
          } catch (err) {
            console.error("Failed to mark notification read:", err);
            isRead = false;
            card.classList.remove("is-read");
            markReadBtn.style.display = "inline-flex";
          }
        },
      },
    },
    ["Mark Read"]
  );

  const children: HTMLElement[] = [leftContent];
  if (!isRead && userId) children.push(markReadBtn);

  const card = createElement("div", {
    class: `notification-card${isRead ? " is-read" : ""}`,
    dataset: { notifCard: "true" },
  }, children);

  return card;
}

function createSystemLogCard(log: SystemLog, onChange?: () => void): HTMLElement {
  let isRead = Boolean(log.isRead);
  const isError = String(log.type || "info").toLowerCase() === "error";
  const isSuccess = String(log.type || "info").toLowerCase() === "success";

  const badge = createElement("span", {
    class: `notification-badge${isError ? " is-error" : isSuccess ? " is-success" : ""}`,
  }, [String(log.type || "info").toUpperCase()]);

  const content = createElement("div", { class: "notification-card__content" }, [
    badge,
    createElement("strong", { class: `notification-card__title${isRead ? " is-read" : ""}` }, [log.title || "System Message"]),
    createElement("p", { class: "notification-card__message" }, [log.message || ""]),
    createElement("small", { class: "notification-card__time" }, [timeAgo(log.createdAt || Date.now())]),
  ]);

  const markReadBtn = createElement(
    "button",
    {
      class: "notification-action-btn",
      events: {
        click: async (e: Event) => {
          const mouseEvent = e as MouseEvent;
          mouseEvent.stopPropagation();
          try {
            const updatedLog = { ...log, isRead: true };
            const saveMethod = (idxDB as any).update || (idxDB as any).put;

            if (saveMethod) {
              await saveMethod(updatedLog);
            }

            isRead = true;
            card.classList.add("is-read");
            markReadBtn.style.display = "none";

            if (onChange) onChange();
          } catch (err) {
            console.error("Failed to update log state in IndexedDB:", err);
          }
        },
      },
    },
    ["Mark Read"]
  );

  const children: HTMLElement[] = [content];
  if (!isRead) children.push(markReadBtn);

  const card = createElement("div", {
    class: `notification-card${isRead ? " is-read" : ""}${isError ? " is-error" : ""}${isSuccess ? " is-success" : ""}`,
  }, children);

  return card;
}

function createSystemActionBar(logs: SystemLog[], onRefresh: () => void): HTMLElement | null {
  if (!logs.length) return null;

  const actionBarChildren: HTMLElement[] = [];
  const unreadExist = logs.some((log) => !log.isRead);

  if (unreadExist) {
    const markAllBtn = createElement(
      "button",
      {
        class: "notification-action-bar__btn notification-action-bar__btn--success",
        events: {
          click: async () => {
            try {
              const saveMethod = (idxDB as any).update || (idxDB as any).put;
              const updatePromises = logs
                .filter((log) => !log.isRead)
                .map((log) => saveMethod({ ...log, isRead: true }));

              await Promise.all(updatePromises);
              onRefresh();
            } catch (error) {
              console.error("Error batch updating logs in IndexedDB:", error);
            }
          },
        },
      },
      ["Mark All Read"]
    );

    actionBarChildren.push(markAllBtn);
  }

  const clearLogsBtn = createElement(
    "button",
    {
      class: "notification-action-bar__btn notification-action-bar__btn--danger",
      events: {
        click: async () => {
          if (!confirm("Clear all stored system logs?")) return;
          try {
            await idxDB.clear();
            onRefresh();
          } catch (error) {
            console.error("Error clearing IndexedDB logs:", error);
          }
        },
      },
    },
    ["Clear All Logs"]
  );

  actionBarChildren.push(clearLogsBtn);

  return createElement("div", { class: "notification-action-bar" }, actionBarChildren);
}

function createActionBar(userId: string | null, notifications: NotificationItem[], onRefresh: () => void): HTMLElement | null {
  if (!userId || !notifications.length) return null;

  const actionBarChildren: HTMLElement[] = [];
  const unreadExist = notifications.some((n) => !n.isRead);

  if (unreadExist) {
    const markAllBtn = createElement(
      "button",
      {
        class: "notification-action-bar__btn notification-action-bar__btn--success",
        events: {
          click: async () => {
            try {
              await markAllNotificationsAsRead();
              onRefresh();
            } catch (error) {
              console.error("Error marking all notifications as read:", error);
            }
          },
        },
      },
      ["Mark All Read"]
    );

    actionBarChildren.push(markAllBtn);
  }

  const clearBtn = createElement(
    "button",
    {
      class: "notification-action-bar__btn notification-action-bar__btn--danger",
      events: {
        click: async () => {
          if (!confirm("Clear all activity notifications?")) return;
          try {
            await clearAllNotifications();
            onRefresh();
          } catch (error) {
            console.error("Error clearing activity notifications:", error);
          }
        },
      },
    },
    ["Clear All"]
  );

  actionBarChildren.push(clearBtn);

  return createElement("div", { class: "notification-action-bar" }, actionBarChildren);
}