// displayCrops.ts
import { createElement } from "../../components/createElement.js";
import { createListOrGrid } from "../../components/layout/listOrGrid.js";

const MOCK_CROPS = Array.from({ length: 12 }, (_, i) => ({
  id: i + 1,
  title: `Crop Variety #${i + 1} (${i % 2 === 0 ? "Wheat" : "Paddy"})`,
  category: i % 2 === 0 ? "Rabi" : "Kharif",
  duration: 110 + (i * 5),
  yieldPotential: `${35 + i * 2} Quintals/Acre`
}));

export function displayCrops(contentContainer: HTMLElement): void {
  contentContainer.replaceChildren();

  let currentPage = 1;
  const pageSize = 10;

  const wrapper = createElement("div", { class: "catalog-page" });
  const header = createElement("div", { class: "catalog-header" });
  header.appendChild(createElement("h2", {}, ["Crop Catalog"]));
  header.appendChild(createElement("p", { class: "text-muted" }, ["Explore certified seed varieties and crop profiles."]));
  wrapper.appendChild(header);

  const contentArea = createElement("div", {});
  wrapper.appendChild(contentArea);

  const paginationBar = createElement("div", { class: "pagination-bar" });
  const prevBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["Previous"]);
  const pageIndicator = createElement("span", { class: "page-indicator" }, ["Page 1"]);
  const nextBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["Next"]);
  paginationBar.append(prevBtn, pageIndicator, nextBtn);
  wrapper.appendChild(paginationBar);

  const renderPage = (page: number) => {
    const start = (page - 1) * pageSize;
    const slice = MOCK_CROPS.slice(start, start + pageSize);

    pageIndicator.textContent = `Page ${page} of ${Math.ceil(MOCK_CROPS.length / pageSize)}`;
    prevBtn.disabled = page === 1;
    nextBtn.disabled = start + pageSize >= MOCK_CROPS.length;

    contentArea.replaceChildren();

    const grid = createListOrGrid({
      items: slice,
      defaultView: "grid",
      renderItem: (crop) => {
        const card = createElement("div", { class: "recipe-card" });
        card.innerHTML = `
          <h4>${crop.title}</h4>
          <div class="card-meta">
            <span class="badge">${crop.category}</span>
            <span class="badge">⏳ ${crop.duration} days</span>
            <span class="badge">📈 ${crop.yieldPotential}</span>
          </div>
        `;
        const viewBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["View Crop Details"]);
        viewBtn.addEventListener("click", () => {
          import("./cropPage.js").then(({ displayCrop }) => {
            displayCrop(contentContainer, true, String(crop.id));
          });
        });
        card.appendChild(viewBtn);
        return card;
      }
    });

    contentArea.appendChild(grid);
  };

  prevBtn.addEventListener("click", () => {
    if (currentPage > 1) {
      currentPage--;
      renderPage(currentPage);
    }
  });

  nextBtn.addEventListener("click", () => {
    if (currentPage * pageSize < MOCK_CROPS.length) {
      currentPage++;
      renderPage(currentPage);
    }
  });

  contentContainer.appendChild(wrapper);
  renderPage(currentPage);
}