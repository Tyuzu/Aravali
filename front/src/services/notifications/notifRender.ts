import { createElement } from "../../components/createElement.js";
import * as idxDB from "../../utils/idxDB.js";
import {
  clearAllNotifications,
  markAllNotificationsAsRead,
  markNotificationAsRead,
  NotificationItem,
} from "./notifService.js";
import { decrementUnreadNotificationState, syncUnreadNotificationState } from "./notifState.js";

export interface SystemLog {
  id?: string | number;
  notificationid?: string | number;
  type?: "info" | "error" | "success" | string;
  title?: string;
  message?: string;
  createdAt?: string | number | Date;
  isRead?: boolean;
}

export type NotificationFilter = "all" | "unread";
export type SystemFilter = "all" | "unread" | "error" | "success";

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

export function filterNotifications(items: NotificationItem[], filter: NotificationFilter, query: string): NotificationItem[] {
  return items.filter((item) => {
    const matchesFilter = filter === "all" ? true : !item.isRead;
    const text = [item.title, item.type, item.message].join(" ");
    return matchesFilter && matchesQuery(text, query);
  });
}

export function filterSystemLogs(items: SystemLog[], filter: SystemFilter, query: string): SystemLog[] {
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

export function renderSummaryChips(
  container: HTMLElement,
  items: Array<{ label: string; value: number; tone?: "default" | "success" | "error" | "warning" }>
): void {
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

export function renderEmptyState(container: HTMLElement, message: string): void {
  container.appendChild(
    createElement("div", { class: "notification-empty-state" }, [
      createElement("p", { class: "notification-empty-state__title" }, ["🔔 All caught up"]),
      createElement("p", { class: "notification-empty-state__message" }, [message]),
    ])
  );
}

export function createNotificationCard(n: NotificationItem, userId: string | null, onChange?: () => void): HTMLElement {
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
            decrementUnreadNotificationState();
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

  const card = createElement(
    "div",
    {
      class: `notification-card${isRead ? " is-read" : ""}`,
      dataset: { notifCard: "true" },
    },
    children
  );

  return card;
}

export function createSystemLogCard(log: SystemLog, onChange?: () => void): HTMLElement {
  let isRead = Boolean(log.isRead);
  const isError = String(log.type || "info").toLowerCase() === "error";
  const isSuccess = String(log.type || "info").toLowerCase() === "success";

  const badge = createElement(
    "span",
    {
      class: `notification-badge${isError ? " is-error" : isSuccess ? " is-success" : ""}`,
    },
    [String(log.type || "info").toUpperCase()]
  );

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

  const card = createElement(
    "div",
    {
      class: `notification-card${isRead ? " is-read" : ""}${isError ? " is-error" : ""}${isSuccess ? " is-success" : ""}`,
    },
    children
  );

  return card;
}

export function createSystemActionBar(logs: SystemLog[], onRefresh: () => void): HTMLElement | null {
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
              syncUnreadNotificationState(0);
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

export function createActionBar(
  userId: string | null,
  notifications: NotificationItem[],
  onRefresh: () => void
): HTMLElement | null {
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
              syncUnreadNotificationState(0);
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
            syncUnreadNotificationState(0);
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
