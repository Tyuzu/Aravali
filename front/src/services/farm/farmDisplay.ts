// displayFarm.ts
export function displayFarm(isLoggedIn: boolean, farmId: string, contentContainer: HTMLElement): void {
  contentContainer.replaceChildren();

  const wrapper = document.createElement("div");
  wrapper.className = "item-detail-container";

  const header = document.createElement("div");
  header.className = "detail-header";
  header.innerHTML = `
    <h1>Green Valley Organic Farm #${farmId}</h1>
    <div class="card-meta" style="margin-top: 8px;">
      <span class="badge">Certified Organic</span>
      <span class="badge">Sonoma County, CA</span>
      <span class="badge">50 Acres</span>
    </div>
  `;

  const desc = document.createElement("p");
  desc.className = "recipe-description";
  desc.textContent = "Green Valley Organic Farm has been family-owned since 1985, specializing in regenerative agriculture, heirloom vegetables, and farm-fresh produce.";

  const practicesSection = document.createElement("div");
  practicesSection.className = "recipe-section";
  practicesSection.innerHTML = `
    <h3>Farming Practices</h3>
    <ul class="ingredient-bullet-list">
      <li>Regenerative Soil Management</li>
      <li>Pesticide-Free Crop Rotation</li>
      <li>Water-Efficient Drip Irrigation</li>
    </ul>
  `;

  const actions = document.createElement("div");
  actions.className = "detail-actions";
  const backBtn = document.createElement("button");
  backBtn.className = "btn-secondary";
  backBtn.textContent = "← Back to Farms";
  backBtn.onclick = () => {
    import("./displayFarms.js").then(({ displayFarms }) => {
      displayFarms(contentContainer, isLoggedIn);
    });
  };
  actions.appendChild(backBtn);

  if (isLoggedIn) {
    const contactBtn = document.createElement("button");
    contactBtn.className = "btn-submit";
    contactBtn.textContent = "Contact Farm Owner";
    contactBtn.onclick = () => alert("Inquiry sent to farm manager!");
    actions.appendChild(contactBtn);
  }

  wrapper.append(header, desc, practicesSection, actions);
  contentContainer.appendChild(wrapper);
}