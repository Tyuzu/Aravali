import "../../../css/ui/Notify.css";
import { createElement } from "../../components/createElement.js"; // Adjust path as needed
import { getState, setState } from "../../state/state.js";
// import { playSoundAlert } from "../../services/notifications/soundAlerts.js";
import { addSystemLog } from "../../utils/idxDB.js"; // Import IndexedDB persistence helper

// ---- Types & Interfaces ----

export type NotificationType = "info" | "success" | "error" | "warning";

export interface NotifyOptions {
  userId?: string | number;
  title?: string;
  type?: NotificationType;
  duration?: number;
  dismissible?: boolean;
  trackUnread?: boolean;
}

function normalizeNotifyOptions(
  options?: NotifyOptions | NotificationType
): NotifyOptions & { type: NotificationType } {
  if (typeof options === "string") {
    return { title: "", type: options as NotificationType, duration: 0, dismissible: true, trackUnread: false };
  }
  return {
    title: "",
    duration: 0,
    dismissible: true,
    trackUnread: false,
    ...options,
    type: options && typeof options === "object" && options.type ? options.type : "info"
  };
}

/**
 * Creates and renders a floating notification alert node.
 */
const Notify = (
  message: string,
  options?: NotifyOptions | NotificationType
): HTMLDivElement => {
  const normalizedOptions = normalizeNotifyOptions(options);
  const { title = "", type = "info", duration = 0, dismissible = true, trackUnread = false } = normalizedOptions;

  const currentUser = getState("user") as { id?: string | number; userid?: string | number } | undefined;
  const currentUserId = normalizedOptions.userId ?? currentUser?.id ?? currentUser?.userid ?? (getState("userid") as string | number | undefined);

  // Track browser window timeouts
  let hideTimeoutId: number | null = null;
  let removeTimeoutId: number | null = null;

  const removeNotification = (): void => {
    if (hideTimeoutId !== null) window.clearTimeout(hideTimeoutId);
    if (removeTimeoutId !== null) window.clearTimeout(removeTimeoutId);
    notify.remove();
  };

  const messageNode = createElement("span", { class: "notify-message" }, [message]);
  const children: HTMLElement[] = [messageNode];

  if (dismissible) {
    const closeBtn = createElement(
      "button",
      {
        class: "notify-close",
        "aria-label": "Close",
        events: {
          click: removeNotification,
        },
      },
      ["×"]
    );
    children.push(closeBtn);
  }

  const notify = createElement(
    "div",
    {
      class: `notify ${type}`,
      role: "alert",
      "aria-live": "assertive",
    },
    children
  ) as HTMLDivElement;

  const containerId = "notify-container";
  let container = document.getElementById(containerId) as HTMLDivElement | null;
  if (!container) {
    container = createElement("div", {
      id: containerId,
      class: "notify-container",
    }) as HTMLDivElement;

    // Fall back safely to document.body if '#app' isn't in the DOM yet
    const appRoot = document.getElementById("app") || document.body;
    appRoot.appendChild(container);
  }

  container.appendChild(notify);

  // Auto-dismiss logic using window.setTimeout
  const timeout = duration || Math.max(3000, message.length * 50);
  hideTimeoutId = window.setTimeout(() => {
    notify.classList.add("hide");
    removeTimeoutId = window.setTimeout(() => notify.remove(), 500);
  }, timeout);

  // Global app state & Side Effects
  if (trackUnread && currentUserId) {
    const unreadTotal = Number(getState("unreadNotifications") || 0);
    const safeUnreadTotal = Number.isFinite(unreadTotal) ? unreadTotal : 0;
    setState("unreadNotifications", safeUnreadTotal + 1);
  }
  // playSoundAlert({ type: "notification" });

  // Persistent System Log Storage (Async save to IndexedDB)
  if (currentUserId) {
    addSystemLog({
      userId: currentUserId,
      title: title || type.charAt(0).toUpperCase() + type.slice(1) + " Alert",
      message: message,
      type: type,
    }).catch((err: unknown) => {
      console.error("Failed to persist notification to IndexedDB:", err);
    });
  } else {
    console.warn("Notification rendered without active user session; skipped IndexedDB persistence.");
  }

  return notify;
};

export default Notify;