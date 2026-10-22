// displayCrop.ts
import { createElement } from "../../components/createElement.js";

export function displayCrop(contentContainer: HTMLElement, isLoggedIn: boolean, cropId: string): void {
  contentContainer.replaceChildren();

  const wrapper = createElement("div", { class: "item-detail-container" });

  const header = createElement("div", { class: "detail-header" });
  header.appendChild(createElement("h1", {}, [`Wheat Variety #${cropId}`]));
  header.appendChild(
    createElement("div", { class: "card-meta", style: "margin-top: 8px;" }, [
      createElement("span", { class: "badge" }, ["Rabi Season"]),
      createElement("span", { class: "badge" }, ["Duration: 120 Days"]),
      createElement("span", { class: "badge" }, ["High Yield"])
    ])
  );
  wrapper.appendChild(header);

  const desc = createElement("p", { class: "recipe-description" }, [
    "High-yielding disease-resistant wheat variety engineered for optimum flour quality and strong gluten content."
  ]);
  wrapper.appendChild(desc);

  const actions = createElement("div", { class: "detail-actions" });
  const aboutBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["View Agronomy Guide"]);
  aboutBtn.addEventListener("click", () => {
    import("./displayAboutCrop.js").then(({ displayAboutCrop }) => {
      displayAboutCrop(contentContainer, isLoggedIn, cropId);
    });
  });

  actions.appendChild(aboutBtn);

  if (isLoggedIn) {
    const saveBtn = createElement("button", { type: "button", class: "btn-submit" }, ["Add to Farm Rotation"]);
    saveBtn.addEventListener("click", () => alert("Crop added to your farm rotation plan!"));
    actions.appendChild(saveBtn);
  }

  wrapper.appendChild(actions);
  contentContainer.appendChild(wrapper);
}