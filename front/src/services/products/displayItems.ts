// displayItems.ts
import { createElement } from "../../components/createElement.js";
import { createListOrGrid } from "../../components/layout/listOrGrid.js";

// Mock database containing different item types
const MOCK_ITEMS_DATABASE: Record<string, any[]> = {
  products: [
    { id: 1, title: "Handmade Ceramic Mug", category: "Handmade", price: 24.99, rating: 4.8 },
    { id: 2, title: "Organic Cotton T-Shirt", category: "Apparel", price: 29.99, rating: 4.6 }
  ],
  tools: [
    { id: 101, title: "Professional Cordless Drill", category: "Power Tools", price: 120.00, condition: "Like New" },
    { id: 102, title: "Heavy Duty Rotary Hammer", category: "Construction", price: 250.00, condition: "Good" }
  ],
  recipes: [
    { id: 201, title: "Classic Margherita Pizza", category: "Main Course", cookTime: 15, rating: 4.9 },
    { id: 202, title: "Avocado Toast with Poached Egg", category: "Breakfast", cookTime: 10, rating: 4.7 }
  ],
  farms: [
    { id: 301, title: "Green Valley Organic Farm", category: "Crop & Produce", acreage: 50, location: "Sonoma County, CA" },
    { id: 302, title: "Sunny Meadows Dairy", category: "Livestock & Dairy", acreage: 120, location: "Vermont" }
  ]
};

export function displayItems(itemType: string, container: HTMLElement, isLoggedIn: boolean): void {
  container.replaceChildren();

  const normalizedType = itemType.toLowerCase();
  const items = MOCK_ITEMS_DATABASE[normalizedType] || [];

  const wrapper = createElement("div", { class: "catalog-page" });
  
  // Header Section
  const header = createElement("div", { class: "catalog-header" });
  header.appendChild(createElement("h2", {}, [`Explore ${itemType.charAt(0).toUpperCase() + itemType.slice(1)}`]));
  header.appendChild(createElement("p", { class: "text-muted" }, [`Browse our community-listed ${normalizedType}.`]));
  wrapper.appendChild(header);

  const dynamicContentArea = createElement("div", { class: "catalog-content" });
  wrapper.appendChild(dynamicContentArea);

  // Use the reusable listOrGrid component
  const listOrGridEl = createListOrGrid({
    items,
    defaultView: "grid",
    renderItem: (item) => {
      const card = createElement("div", { class: "item-card" });
      const title = createElement("h4", {}, [item.title]);
      
      const meta = createElement("div", { class: "card-meta" }, [
        createElement("span", { class: "badge" }, [item.category]),
        item.price ? createElement("span", {}, [`$${item.price.toFixed(2)}`]) : null,
        item.cookTime ? createElement("span", {}, [`⏱ ${item.cookTime}m`]) : null,
        item.acreage ? createElement("span", {}, [`🌾 ${item.acreage} acres`]) : null
      ].filter(Boolean) as HTMLElement[]);

      const viewBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["View Details"]);
      viewBtn.addEventListener("click", () => {
        // Trigger single item detail view
        import("./productPage.js").then(({ displayProduct }) => {
          displayProduct(container, isLoggedIn, item.id, normalizedType);
        });
      });

      card.append(title, meta, viewBtn);
      return card;
    }
  });

  dynamicContentArea.appendChild(listOrGridEl);
  container.appendChild(wrapper);
}