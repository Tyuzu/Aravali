// displayRecipes.ts
import { createElement } from "../../components/createElement.js";
import { createListOrGrid } from "../../components/layout/listOrGrid.js";

interface RecipeSummary {
  id: number;
  title: string;
  category: string;
  cookTime: number;
  rating: number;
}

// Generate 10 mock recipes for testing pagination / infinite scroll
const ALL_RECIPES: RecipeSummary[] = Array.from({ length: 25 }, (_, i) => ({
  id: i + 1,
  title: `Delicious Recipe #${i + 1}`,
  category: i % 2 === 0 ? "Main Course" : "Dessert",
  cookTime: 15 + (i * 5) % 45,
  rating: Number((4.0 + (i % 10) * 0.1).toFixed(1))
}));

export function displayRecipes(container: HTMLElement, isLoggedIn: boolean): void {
  container.replaceChildren();

  let currentPage = 1;
  const pageSize = 10;

  const wrapper = createElement("div", { class: "recipes-catalog-page" });
  const headingSection = createElement("div", { class: "catalog-header" });
  headingSection.appendChild(createElement("h2", {}, ["Explore Recipes"]));
  headingSection.appendChild(createElement("p", { class: "text-muted" }, ["Discover hand-crafted recipes from our community."]));
  wrapper.appendChild(headingSection);

  const dynamicContentArea = createElement("div", { class: "catalog-content" });
  wrapper.appendChild(dynamicContentArea);

  // Pagination Controls Bar
  const paginationBar = createElement("div", { class: "pagination-bar" });
  const prevBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["Previous"]);
  const pageIndicator = createElement("span", { class: "page-indicator" }, ["Page 1"]);
  const nextBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["Next"]);
  paginationBar.append(prevBtn, pageIndicator, nextBtn);
  wrapper.appendChild(paginationBar);

  const renderPage = (page: number) => {
    const startIndex = (page - 1) * pageSize;
    const currentItems = ALL_RECIPES.slice(startIndex, startIndex + pageSize);

    pageIndicator.textContent = `Page ${page} of ${Math.ceil(ALL_RECIPES.length / pageSize)}`;
    prevBtn.disabled = page === 1;
    nextBtn.disabled = startIndex + pageSize >= ALL_RECIPES.length;

    dynamicContentArea.replaceChildren();

    // Use the reusable listOrGrid component here
    const listOrGridEl = createListOrGrid({
      items: currentItems,
      defaultView: "grid",
      renderItem: (recipe) => {
        const card = createElement("div", { class: "recipe-card" });
        const cardTitle = createElement("h4", {}, [recipe.title]);
        const cardMeta = createElement("div", { class: "card-meta" }, [
          createElement("span", {}, [recipe.category]),
          createElement("span", {}, [`⏱ ${recipe.cookTime}m`]),
          createElement("span", {}, [`★ ${recipe.rating}`])
        ]);

        const viewBtn = createElement("button", { type: "button", class: "btn-secondary" }, ["View Recipe"]);
        viewBtn.addEventListener("click", () => {
          // Trigger single recipe view
          import("./recipePage.js").then(({ displayRecipe }) => {
            displayRecipe(container, isLoggedIn, recipe.id);
          });
        });

        card.append(cardTitle, cardMeta, viewBtn);
        return card;
      }
    });

    dynamicContentArea.appendChild(listOrGridEl);
  };

  prevBtn.addEventListener("click", () => {
    if (currentPage > 1) {
      currentPage--;
      renderPage(currentPage);
    }
  });

  nextBtn.addEventListener("click", () => {
    if (currentPage * pageSize < ALL_RECIPES.length) {
      currentPage++;
      renderPage(currentPage);
    }
  });

  container.appendChild(wrapper);
  renderPage(currentPage);
}