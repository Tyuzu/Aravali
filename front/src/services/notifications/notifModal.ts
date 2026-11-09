import "../../../css/subpages/notifications.css";
import Modal from "../../components/ui/Modal.js";
import { createElement } from "../../components/createElement.js";
import { getNotifications, NotificationItem } from "./notifService.js";
import * as idxDB from "../../utils/idxDB.js";
import { getUserId } from "../../utils/getUserID.js";
import {
  createActionBar,
  createNotificationCard,
  createSystemActionBar,
  createSystemLogCard,
  filterNotifications,
  filterSystemLogs,
  renderEmptyState,
  renderSummaryChips,
  SummaryChip,
  SystemLog
} from "./notifRender.js";

const UI_STATE = {
  search: "",
  activityFilter: "all",
  systemFilter: "all",
};

export async function openNotificationsModal(): Promise<void> {
  const userId = getUserId();
  let activeTab: "activity" | "system" = "activity";

  const content = createElement("div", { class: "notification-modal" });
  const tabHeader = createElement("div", { class: "notification-tab-header" });
  const activityTabBtn = createElement("button", { class: "notification-tab-button" }, ["Activity"]);
  const systemTabBtn = createElement("button", { class: "notification-tab-button is-active" }, ["System Logs"]);

  const toolbar = createElement("div", { class: "notification-toolbar" });
  const searchInput = createElement("input", {
    class: "notification-search",
    type: "search",
    placeholder: "Search activity...",
    value: UI_STATE.search,
  }) as HTMLInputElement;

  const filterGroup = createElement("div", { class: "notification-filter-group" });
  const activityFilterButtons: Record<string, HTMLButtonElement> = {
    all: createElement("button", { class: "notification-filter-btn is-active", type: "button" }, ["All"]) as HTMLButtonElement,
    unread: createElement("button", { class: "notification-filter-btn", type: "button" }, ["Unread"]) as HTMLButtonElement,
  };

  function renderActivityFilters(): void {
    filterGroup.replaceChildren();
    Object.entries(activityFilterButtons).forEach(([key, button]) => {
      button.classList.toggle("is-active", UI_STATE.activityFilter === key);
      button.onclick = () => {
        UI_STATE.activityFilter = key;
        renderActivityFilters();
        renderActivityTab();
      };
      filterGroup.appendChild(button);
    });
  }

  function renderSystemFilters(): void {
    const systemFilterButtons: Record<string, HTMLButtonElement> = {
      all: createElement("button", { class: "notification-filter-btn is-active", type: "button" }, ["All"]) as HTMLButtonElement,
      unread: createElement("button", { class: "notification-filter-btn", type: "button" }, ["Unread"]) as HTMLButtonElement,
      error: createElement("button", { class: "notification-filter-btn", type: "button" }, ["Errors"]) as HTMLButtonElement,
      success: createElement("button", { class: "notification-filter-btn", type: "button" }, ["Success"]) as HTMLButtonElement,
    };

    filterGroup.replaceChildren();
    Object.entries(systemFilterButtons).forEach(([key, button]) => {
      button.classList.toggle("is-active", UI_STATE.systemFilter === key);
      button.onclick = () => {
        UI_STATE.systemFilter = key;
        renderSystemFilters();
        renderSystemTab();
      };
      filterGroup.appendChild(button);
    });
  }

  toolbar.appendChild(searchInput);
  toolbar.appendChild(filterGroup);
  tabHeader.appendChild(systemTabBtn);
  tabHeader.appendChild(activityTabBtn);
  content.appendChild(tabHeader);
  content.appendChild(toolbar);

  const tabContentView = createElement("div", { class: "notification-tab-content" });
  content.appendChild(tabContentView);

  Modal({ title: "📬 Notifications & Logs", content, size: "medium", showCloseButton: true });

  searchInput.addEventListener("input", (event: Event) => {
    const target = event.target as HTMLInputElement | null;
    UI_STATE.search = target ? target.value : "";
    if (activeTab === "activity") renderActivityTab();
    else renderSystemTab();
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
    searchInput.placeholder = "Search activity...";
    renderActivityFilters();
  }

  function applySystemFilterState(): void {
    searchInput.placeholder = "Search system logs...";
    renderSystemFilters();
  }

  function updateTabStyles(): void {
    const isActivity = activeTab === "activity";
    activityTabBtn.classList.toggle("is-active", isActivity);
    systemTabBtn.classList.toggle("is-active", !isActivity);
  }

  async function renderActivityTab(): Promise<void> {
    tabContentView.innerHTML = "";
    const summaryHost = createElement("div");
    tabContentView.appendChild(summaryHost);
    summaryHost.innerHTML = '<div class="notification-loading">Loading activity...</div>';

    try {
      const notifications: NotificationItem[] = (await getNotifications()) || [];
      const filtered = filterNotifications(notifications, UI_STATE.activityFilter, UI_STATE.search);
      const unreadCount = notifications.filter((n) => !n.isRead).length;

      const summaryChips: SummaryChip[] = [
        { label: "Total", value: notifications.length },
        { label: "Unread", value: unreadCount, tone: "warning" },
      ];
      renderSummaryChips(summaryHost, summaryChips);

      const actionBar = createActionBar(userId, filtered, renderActivityTab);
      if (actionBar) tabContentView.appendChild(actionBar);

      if (!filtered.length) {
        renderEmptyState(tabContentView, "No matching activity updates yet.");
        return;
      }

      const listContainer = createElement("div", { class: "notification-list" });
      filtered.forEach((notification: NotificationItem) => {
        listContainer.appendChild(createNotificationCard(notification, userId, renderActivityTab));
      });
      tabContentView.appendChild(listContainer);
    } catch (err) {
      console.error("Failed to load activity notifications:", err);
      tabContentView.innerHTML = '<div class="notification-error">Failed to load activity notifications.</div>';
    }
  }

  async function renderSystemTab(): Promise<void> {
    tabContentView.innerHTML = "";
    const summaryHost = createElement("div");
    tabContentView.appendChild(summaryHost);
    summaryHost.innerHTML = '<div class="notification-loading">Loading system logs...</div>';

    let logs: SystemLog[] = [];
    try {
      logs = userId ? ((await idxDB.getAll(userId)) || []) : [];
      logs.sort((a, b) => new Date(b.createdAt || 0).getTime() - new Date(a.createdAt || 0).getTime());
    } catch (err) {
      console.error("Failed to fetch system logs from IndexedDB:", err);
    }

    const filtered = filterSystemLogs(logs, UI_STATE.systemFilter, UI_STATE.search);
    const unreadCount = logs.filter((log) => !log.isRead).length;
    const errorCount = logs.filter((log) => String(log.type || "info").toLowerCase() === "error").length;

    const summaryChips: SummaryChip[] = [
      { label: "Total", value: logs.length },
      { label: "Unread", value: unreadCount, tone: "warning" },
      { label: "Errors", value: errorCount, tone: "error" },
    ];
    renderSummaryChips(summaryHost, summaryChips);

    if (!filtered.length) {
      renderEmptyState(tabContentView, "No matching system logs found.");
      return;
    }

    const actionBar = createSystemActionBar(userId, filtered, renderSystemTab);
    if (actionBar) tabContentView.appendChild(actionBar);

    const listContainer = createElement("div", { class: "notification-list" });
    filtered.forEach((log) => {
      listContainer.appendChild(createSystemLogCard(userId, log, renderSystemTab));
    });
    tabContentView.appendChild(listContainer);
  }

  applyActivityFilterState();
  updateTabStyles();
  await renderActivityTab();
}