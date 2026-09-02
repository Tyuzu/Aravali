import "../../../css/inistyles/adminpage.css";
import { approveRoleRequest, listRoleRequests, rejectRoleRequest, RoleApplication } from "../../services/admin/roleManagement.js";

export async function RoleRequestsPage(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";

  if (!isLoggedIn) {
    contentContainer.innerHTML = `<div class="empty-state">Please log in to view role requests.</div>`;
    return;
  }

  const wrapper = document.createElement("div");
  wrapper.className = "admin-panel";

  const heading = document.createElement("h2");
  heading.textContent = "Role Requests";
  wrapper.appendChild(heading);

  const statusFilter = document.createElement("select");
  statusFilter.innerHTML = `
    <option value="">All</option>
    <option value="pending">Pending</option>
    <option value="approved">Approved</option>
    <option value="rejected">Rejected</option>
  `;

  const listHost = document.createElement("div");
  listHost.className = "admin-list";
  wrapper.appendChild(statusFilter);
  wrapper.appendChild(listHost);
  contentContainer.appendChild(wrapper);

  async function render(status = "") {
    listHost.innerHTML = "Loading...";
    const items = await listRoleRequests(status);

    if (!items || items.length === 0) {
      listHost.innerHTML = "<div class=\"empty-state\">No role requests found.</div>";
      return;
    }

    listHost.innerHTML = "";
    items.forEach((item: RoleApplication) => {
      const card = document.createElement("div");
      card.className = "card";

      const title = document.createElement("h3");
      title.textContent = `${item.role} request`;

      const user = document.createElement("p");
      user.textContent = `User: ${item.user_id}`;

      const reason = document.createElement("p");
      reason.textContent = `Reason: ${item.reason}`;

      const statusPill = document.createElement("span");
      statusPill.textContent = item.status;
      statusPill.style.textTransform = "capitalize";

      const actions = document.createElement("div");
      actions.style.display = "flex";
      actions.style.gap = "8px";
      if (item.status === "pending") {
        const approveBtn = document.createElement("button");
        approveBtn.textContent = "Approve";
        approveBtn.addEventListener("click", async () => {
          await approveRoleRequest(item.id);
          await render(statusFilter.value);
        });

        const rejectBtn = document.createElement("button");
        rejectBtn.textContent = "Reject";
        rejectBtn.addEventListener("click", async () => {
          await rejectRoleRequest(item.id);
          await render(statusFilter.value);
        });

        actions.appendChild(approveBtn);
        actions.appendChild(rejectBtn);
      }

      card.append(title, user, reason, statusPill, actions);
      listHost.appendChild(card);
    });
  }

  statusFilter.addEventListener("change", () => render(statusFilter.value));
  await render();
}
