import { listRoleRequests } from "./roleManagement.js";
import { navigate } from "../../routes/navigate.js";

/**
 * @param {HTMLElement} container
 * @param {boolean} isLoggedIn
 * @returns {Promise<void>}
 */
export async function displayAdminDash(container, isLoggedIn) {
  container.replaceChildren();

  const panel = document.createElement("div");
  panel.className = "admin-panel";

  const heading = document.createElement("h2");
  heading.textContent = "Admin Dashboard";
  panel.appendChild(heading);

  if (!isLoggedIn) {
    const msg = document.createElement("div");
    msg.className = "empty-state";
    msg.textContent = "Please log in to access the admin dashboard.";
    panel.appendChild(msg);
    container.appendChild(panel);
    return;
  }

  const actions = document.createElement("div");
  actions.style.display = "flex";
  actions.style.gap = "12px";
  actions.style.flexWrap = "wrap";
  actions.style.marginBottom = "16px";

  const roleRequestBtn = document.createElement("button");
  roleRequestBtn.type = "button";
  roleRequestBtn.textContent = "Review Role Requests";
  roleRequestBtn.addEventListener("click", () => navigate("/admin/role-requests"));

  const moderatorBtn = document.createElement("button");
  moderatorBtn.type = "button";
  moderatorBtn.textContent = "Review Moderator Applications";
  moderatorBtn.addEventListener("click", () => navigate("/admin/moderator-applications"));

  actions.appendChild(roleRequestBtn);
  actions.appendChild(moderatorBtn);
  panel.appendChild(actions);

  const summary = document.createElement("div");
  summary.className = "admin-list";
  summary.style.display = "grid";
  summary.style.gridTemplateColumns = "repeat(auto-fit, minmax(160px, 1fr))";
  summary.style.gap = "12px";
  panel.appendChild(summary);

  try {
    const requests = await listRoleRequests();
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
    const errorBox = document.createElement("div");
    errorBox.className = "empty-state";
    errorBox.textContent = "Unable to load admin summary right now.";
    summary.appendChild(errorBox);
    console.error("Failed to load admin dashboard summary", error);
  }

  container.appendChild(panel);
}
