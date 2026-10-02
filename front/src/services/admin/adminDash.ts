import { listRoleRequests, canReviewRoleRequests, normalizeRoleName, type RoleApplication } from "./roleManagement.js";
import { navigate } from "../../routes/navigate.js";
import { getState } from "../../state/state.js";

export async function displayAdminDash(container: HTMLElement, isLoggedIn: boolean): Promise<void> {
  container.replaceChildren();

  const panel: HTMLDivElement = document.createElement("div");
  panel.className = "admin-panel";

  const heading: HTMLHeadingElement = document.createElement("h2");
  heading.textContent = "Admin Dashboard";
  panel.appendChild(heading);

  if (!isLoggedIn) {
    const msg: HTMLDivElement = document.createElement("div");
    msg.className = "empty-state";
    msg.textContent = "Please log in to access the admin dashboard.";
    panel.appendChild(msg);
    container.appendChild(panel);
    return;
  }

  const rawRoles: unknown = getState("roles");
  const userRoles: string[] = Array.isArray(rawRoles)
    ? rawRoles.map((role: unknown) => normalizeRoleName(role))
    : [];
  const canReviewRoleRequestsUI = canReviewRoleRequests(userRoles);
  const canReviewModeratorApps = userRoles.includes("admin");

  const actions: HTMLDivElement = document.createElement("div");
  actions.style.display = "flex";
  actions.style.gap = "12px";
  actions.style.flexWrap = "wrap";
  actions.style.marginBottom = "16px";

  if (canReviewRoleRequestsUI) {
    const roleRequestBtn: HTMLButtonElement = document.createElement("button");
    roleRequestBtn.type = "button";
    roleRequestBtn.textContent = "Review Role Requests";
    roleRequestBtn.addEventListener("click", () => {
      void navigate("/admin/role-requests");
    });
    actions.appendChild(roleRequestBtn);
  }

  if (canReviewModeratorApps) {
    const moderatorBtn: HTMLButtonElement = document.createElement("button");
    moderatorBtn.type = "button";
    moderatorBtn.textContent = "Review Moderator Applications";
    moderatorBtn.addEventListener("click", () => {
      void navigate("/admin/moderator-applications");
    });
    actions.appendChild(moderatorBtn);
  }

  if (actions.children.length > 0) {
    panel.appendChild(actions);
  }

  if (actions.children.length === 0) {
    const msg: HTMLDivElement = document.createElement("div");
    msg.className = "empty-state";
    msg.textContent = "You do not have access to the admin dashboard.";
    panel.appendChild(msg);
    container.appendChild(panel);
    return;
  }

  const summary: HTMLDivElement = document.createElement("div");
  summary.className = "admin-list";
  summary.style.display = "grid";
  summary.style.gridTemplateColumns = "repeat(auto-fit, minmax(160px, 1fr))";
  summary.style.gap = "12px";
  panel.appendChild(summary);

  try {
    const requests: RoleApplication[] = await listRoleRequests();
    const pending = requests.filter((item) => item?.status === "pending").length;
    const approved = requests.filter((item) => item?.status === "approved").length;
    const rejected = requests.filter((item) => item?.status === "rejected").length;

    const cards = [
      { label: "Pending", value: pending },
      { label: "Approved", value: approved },
      { label: "Rejected", value: rejected },
      { label: "Total", value: requests.length }
    ];

    cards.forEach((card) => {
      const item = document.createElement("div");
      item.className = "card";
      item.innerHTML = `<strong>${card.label}</strong><div>${card.value}</div>`;
      summary.appendChild(item);
    });
  } catch (error) {
    const errorBox: HTMLDivElement = document.createElement("div");
    errorBox.className = "empty-state";
    errorBox.textContent = "Unable to load admin summary right now.";
    summary.appendChild(errorBox);
    console.error("Failed to load admin dashboard summary", error);
  }

  container.appendChild(panel);
}
