import "../../../css/subpages/notifications.css";
import Modal from "../../components/ui/Modal.js";
import { createElement } from "../../components/createElement.js";
import { getNotifications } from "./notifService.js";
import * as idxDB from "../../utils/idxDB.js";
import { getUserId, syncUnreadNotificationState } from "./notifState.js";
import {
  createActionBar,
  createNotificationCard,
  createSystemActionBar,
  createSystemLogCard,
  filterNotifications,
  filterSystemLogs,
  renderEmptyState,
  renderSummaryChips,
  type NotificationFilter,
  type SystemFilter,
  type SystemLog,
} from "./notifRender.js";

const UI_STATE = {
  search: "",
  activityFilter: "all" as NotificationFilter,
  systemFilter: "all" as SystemFilter,
};

export async function openNotificationsModal(): Promise<void> {
  const userId = getUserId();
  let activeTab: "activity" | "system" = "activity";

  syncUnreadNotificationState(0);

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

  Modal({ title: "📬 Notifications & Logs", content, size: "medium", showCloseButton: true });

  searchInput.addEventListener("input", (event: Event) => {
    const value = (event.target as HTMLInputElement).value;
    UI_STATE.search = value;
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
    summaryHost.innerHTML = '<div class="notification-loading">Loading activity...</div>';

    try {
      const notifications = (await getNotifications()) || [];
      const filtered = filterNotifications(notifications, UI_STATE.activityFilter, UI_STATE.search);
      const unreadCount = notifications.filter((n) => !n.isRead).length;

      renderSummaryChips(summaryHost, [
        { label: "Total", value: notifications.length },
        { label: "Unread", value: unreadCount, tone: "warning" },
      ]);

      const actionBar = createActionBar(userId, filtered, renderActivityTab);
      if (actionBar) tabContentView.appendChild(actionBar);

      if (!filtered.length) {
        renderEmptyState(tabContentView, "No matching activity updates yet.");
        return;
      }

      const listContainer = createElement("div", { class: "notification-list" });
      filtered.forEach((notification) => {
        listContainer.appendChild(createNotificationCard(notification, userId, renderActivityTab));
      });
      tabContentView.appendChild(listContainer);
    } catch (err) {
      console.error("Failed to load activity notifications:", err);
      tabContentView.innerHTML = '<div class="notification-error">Failed to load activity notifications.</div>';
    }
  }

  async function renderSystemTab(): Promise<void> {
    const summaryHost = createElement("div");
    tabContentView.innerHTML = "";
    tabContentView.appendChild(summaryHost);
    summaryHost.innerHTML = '<div class="notification-loading">Loading system logs...</div>';

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

    const actionBar = createSystemActionBar(filtered, renderSystemTab);
    if (actionBar) tabContentView.appendChild(actionBar);

    const listContainer = createElement("div", { class: "notification-list" });
    filtered.forEach((log) => {
      listContainer.appendChild(createSystemLogCard(log, renderSystemTab));
    });
    tabContentView.appendChild(listContainer);
  }

  applyActivityFilterState();
  updateTabStyles();
  await renderActivityTab();
}
