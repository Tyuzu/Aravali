// displayFarmDash.ts
export function displayFarmDash(contentContainer: HTMLElement, isLoggedIn: boolean): void {
  contentContainer.replaceChildren();

  if (!isLoggedIn) {
    const note = document.createElement("div");
    note.className = "empty-state";
    note.textContent = "Please sign in to access your farm dashboard.";
    contentContainer.appendChild(note);
    return;
  }

  const wrapper = document.createElement("div");
  wrapper.className = "farm-dash-container";

  const header = document.createElement("div");
  header.className = "dash-header";
  header.innerHTML = `
    <h1>My Farm Dashboard</h1>
    <p class="text-muted">Manage your farm profile, inventory, and visitor schedules.</p>
  `;

  // Farm Summary Metrics
  const metrics = document.createElement("div");
  metrics.className = "form-row-three";
  metrics.innerHTML = `
    <div class="recipe-card">
      <h4>Active Listings</h4>
      <p class="text-muted">12 Products & Tools</p>
    </div>
    <div class="recipe-card">
      <h4>Farm Acreage</h4>
      <p class="text-muted">50 Acres (Certified Organic)</p>
    </div>
    <div class="recipe-card">
      <h4>Profile Status</h4>
      <p style="color: var(--color-success); font-weight: var(--font-weight-medium);">Active & Verified</p>
    </div>
  `;

  const actions = document.createElement("div");
  actions.className = "detail-actions";
  const editBtn = document.createElement("button");
  editBtn.className = "btn-submit";
  editBtn.textContent = "Edit Farm Profile";
  editBtn.onclick = () => alert("Opening farm profile editor...");
  actions.appendChild(editBtn);

  wrapper.append(header, metrics, actions);
  contentContainer.appendChild(wrapper);
}