// listOrGrid.ts
import { createElement } from "../createElement.js";

export interface ListOrGridConfig<T> {
  items: T[];
  renderItem: (item: T) => HTMLElement;
  defaultView?: "grid" | "list";
  className?: string;
}

export function createListOrGrid<T>({
  items,
  renderItem,
  defaultView = "grid",
  className = ""
}: ListOrGridConfig<T>): HTMLElement {
  const container = createElement("div", { class: `list-grid-wrapper ${className}` });

  // --- Toolbar with View Switcher ---
  const toolbar = createElement("div", { class: "list-grid-toolbar" });
  const viewToggleGroup = createElement("div", { class: "view-toggle-group" });

  const gridBtn = createElement("button", { type: "button", class: `btn-view-toggle ${defaultView === "grid" ? "active" : ""}` }, ["Grid View"]);
  const listBtn = createElement("button", { type: "button", class: `btn-view-toggle ${defaultView === "list" ? "active" : ""}` }, ["List View"]);

  viewToggleGroup.append(gridBtn, listBtn);
  toolbar.appendChild(viewToggleGroup);

  // --- Items Container ---
  const contentWrapper = createElement("div", { class: `items-container ${defaultView}` });

  const renderItems = () => {
    contentWrapper.replaceChildren();
    if (!items || items.length === 0) {
      const emptyState = createElement("div", { class: "empty-state" }, ["No items available to display."]);
      contentWrapper.appendChild(emptyState);
      return;
    }

    items.forEach(item => {
      const itemEl = renderItem(item);
      contentWrapper.appendChild(itemEl);
    });
  };

  // Event Listeners for Switching Views
  gridBtn.addEventListener("click", () => {
    gridBtn.classList.add("active");
    listBtn.classList.remove("active");
    contentWrapper.className = "items-container grid";
  });

  listBtn.addEventListener("click", () => {
    listBtn.classList.add("active");
    gridBtn.classList.remove("active");
    contentWrapper.className = "items-container list";
  });

  renderItems();
  container.append(toolbar, contentWrapper);
  return container;
}