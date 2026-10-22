// displayProduct.ts
import { createElement } from "../../components/createElement.js";

export function displayProduct(
  container: HTMLElement,
  isLoggedIn: boolean,
  productid: string | number,
  productType: string
): void {
  container.replaceChildren();

  const wrapper = createElement("div", { class: "item-detail-container" });

  // Header Section
  const header = createElement("div", { class: "detail-header" });
  const title = createElement("h1", {}, [`${productType.toUpperCase()} #${productid}`]);
  const subtitle = createElement("p", { class: "text-muted" }, [`Detailed view for ${productType} item ID: ${productid}`]);
  header.append(title, subtitle);

  // Content Specifications
  const specs = createElement("div", { class: "detail-specs" });
  specs.append(
    createElement("p", {}, ["Here you can find all specific attributes, descriptions, pricing, options, and interactive widgets associated with this listing."])
  );

  // Action Buttons
  const actions = createElement("div", { class: "detail-actions" });
  const backBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["← Back to List"]);
  backBtn.addEventListener("click", () => {
    import("./displayItems.js").then(({ displayItems }) => {
      displayItems(productType, container, isLoggedIn);
    });
  });
  actions.appendChild(backBtn);

  if (isLoggedIn) {
    const actionBtn = createElement("button", { type: "button", class: "btn-submit" }, ["Contact / Order Listing"]);
    actionBtn.addEventListener("click", () => alert(`Action processed for ${productType} #${productid}!`));
    actions.appendChild(actionBtn);
  }

  wrapper.append(header, specs, actions);
  container.appendChild(wrapper);
}