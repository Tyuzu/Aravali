// displayFarms.ts
import { createListOrGrid } from "../../components/layout/listOrGrid.js";

const MOCK_FARMS = Array.from({ length: 15 }, (_, i) => ({
  id: i + 1,
  title: `Green Valley Farm #${i + 1}`,
  category: i % 2 === 0 ? "Crop & Produce" : "Livestock & Dairy",
  acreage: 25 + i * 15,
  location: `Region ${String.fromCharCode(65 + (i % 5))}`
}));

export function displayFarms(contentContainer: HTMLElement, isLoggedIn: boolean): void {
  contentContainer.replaceChildren();

  let currentPage = 1;
  const pageSize = 10;

  const wrapper = document.createElement("div");
  wrapper.className = "catalog-page";

  const header = document.createElement("div");
  header.className = "catalog-header";
  header.innerHTML = `
    <h2>Explore Local Farms</h2>
    <p class="text-muted">Discover sustainable farms and agricultural producers in your area.</p>
  `;
  wrapper.appendChild(header);

  const contentArea = document.createElement("div");
  wrapper.appendChild(contentArea);

  const paginationBar = document.createElement("div");
  paginationBar.className = "pagination-bar";
  const prevBtn = document.createElement("button");
  prevBtn.className = "btn-secondary";
  prevBtn.textContent = "Previous";
  const pageIndicator = document.createElement("span");
  pageIndicator.className = "page-indicator";
  const nextBtn = document.createElement("button");
  nextBtn.className = "btn-secondary";
  nextBtn.textContent = "Next";
  paginationBar.append(prevBtn, pageIndicator, nextBtn);
  wrapper.appendChild(paginationBar);

  const renderPage = (page: number) => {
    const start = (page - 1) * pageSize;
    const slice = MOCK_FARMS.slice(start, start + pageSize);

    pageIndicator.textContent = `Page ${page} of ${Math.ceil(MOCK_FARMS.length / pageSize)}`;
    prevBtn.disabled = page === 1;
    nextBtn.disabled = start + pageSize >= MOCK_FARMS.length;

    contentArea.replaceChildren();

    const grid = createListOrGrid({
      items: slice,
      defaultView: "grid",
      renderItem: (farm) => {
        const card = document.createElement("div");
        card.className = "recipe-card";
        card.innerHTML = `
          <h4>${farm.title}</h4>
          <div class="card-meta">
            <span class="badge">${farm.category}</span>
            <span class="badge">🌾 ${farm.acreage} acres</span>
            <span class="badge">📍 ${farm.location}</span>
          </div>
        `;
        const btn = document.createElement("button");
        btn.className = "btn-secondary";
        btn.textContent = "View Farm Profile";
        btn.onclick = () => {
          import("./farmDisplay.js").then(({ displayFarm }) => {
            displayFarm(isLoggedIn, String(farm.id), contentContainer);
          });
        };
        card.appendChild(btn);
        return card;
      }
    });

    contentArea.appendChild(grid);
  };

  prevBtn.onclick = () => {
    if (currentPage > 1) {
      currentPage--;
      renderPage(currentPage);
    }
  };

  nextBtn.onclick = () => {
    if (currentPage * pageSize < MOCK_FARMS.length) {
      currentPage++;
      renderPage(currentPage);
    }
  };

  contentContainer.appendChild(wrapper);
  renderPage(currentPage);
}