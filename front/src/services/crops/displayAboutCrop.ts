// displayAboutCrop.ts
import { createElement } from "../../components/createElement.js";

export function displayAboutCrop(contentContainer: HTMLElement, isLoggedIn: boolean, cropId: string): void {
  contentContainer.replaceChildren();

  const wrapper = createElement("div", { class: "item-detail-container" });

  const header = createElement("div", { class: "detail-header" });
  header.appendChild(createElement("h1", {}, [`Crop Knowledge Base: Crop #${cropId}`]));
  header.appendChild(createElement("p", { class: "text-muted" }, ["Comprehensive agronomic guidelines, soil requirements, and pest management."]));
  wrapper.appendChild(header);

  const content = createElement("div", { class: "recipe-section" });
  content.innerHTML = `
    <h3>Ideal Climate & Soil</h3>
    <p>Requires well-drained loamy soil with a pH range of 6.0 to 7.5. Thrives in moderate temperatures between 20°C and 30°C with adequate sunshine.</p>
    
    <h3>Water Management</h3>
    <p>Critical irrigation stages include crown root initiation, tillering, and grain filling. Avoid waterlogging.</p>
  `;
  wrapper.appendChild(content);

  const backBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["← Back to Crops"]);
  backBtn.addEventListener("click", () => {
    import("./displayCrops.js").then(({ displayCrops }) => {
      displayCrops(contentContainer);
    });
  });
  wrapper.appendChild(backBtn);

  contentContainer.appendChild(wrapper);
}