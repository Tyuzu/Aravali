import "../../../css/inistyles/adminpage.css";
import {
  approveModeratorApplication,
  listModeratorApplications,
  rejectModeratorApplication,
  type ModeratorApplication,
} from "../../services/admin/moderatorManagement.js";

export async function ModeratorApplicationsPage(
  isLoggedIn: boolean,
  contentContainer: HTMLElement
): Promise<void> {
  contentContainer.innerHTML = "";

  if (!isLoggedIn) {
    contentContainer.innerHTML = `<div class="empty-state">Please log in to review moderator applications.</div>`;
    return;
  }

  const wrapper = document.createElement("div");
  wrapper.className = "admin-panel";

  const heading = document.createElement("h2");
  heading.textContent = "Moderator Applications";
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
    const selectedStatus = (status || statusFilter.value || "").trim().toLowerCase();
    listHost.innerHTML = "Loading...";

    try {
      const items = await listModeratorApplications(selectedStatus);

      if (!items || items.length === 0) {
        listHost.innerHTML = "<div class=\"empty-state\">No moderator applications found.</div>";
        return;
      }

      listHost.innerHTML = "";
      items.forEach((item: ModeratorApplication) => {
        const card = document.createElement("div");
        card.className = "card";

        const title = document.createElement("h3");
        title.textContent = `Moderator application`;

        const user = document.createElement("p");
        user.textContent = `User: ${item.user_id}`;

        const reason = document.createElement("p");
        reason.textContent = `Reason: ${item.reason || "No reason provided."}`;

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
            await approveModeratorApplication(item.id);
            await render(statusFilter.value);
          });

          const rejectBtn = document.createElement("button");
          rejectBtn.textContent = "Reject";
          rejectBtn.addEventListener("click", async () => {
            await rejectModeratorApplication(item.id);
            await render(statusFilter.value);
          });

          actions.appendChild(approveBtn);
          actions.appendChild(rejectBtn);
        }

        card.append(title, user, reason, statusPill, actions);
        listHost.appendChild(card);
      });
    } catch (error) {
      const msg = error instanceof Error ? error.message : "Unable to load moderator applications.";
      listHost.innerHTML = `<div class="empty-state">${msg}</div>`;
    }
  }

  statusFilter.addEventListener("change", () => render(statusFilter.value));
  await render();
}
