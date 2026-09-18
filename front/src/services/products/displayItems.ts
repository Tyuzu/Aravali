import { createElement } from "../../components/createElement.js";
import Button from "../../components/base/Button.js";
import { fetchFarmItems } from "./api.js";
import { openItemFormModal } from "./createOrEdit.js";
import { renderItemCard } from "./renderItemCard.js";
import { renderCategoryChips } from "./renderCategoryChips.js";
import { capitalize } from "../profile/profileHelpers.js";
import { renderSearchAndSortUI } from "./renderSearchAndSortUI.js";
import { renderPagination } from "./renderPagination.js";
import { DisplayItemsOptions, FarmItem, ItemType } from "./types.js";

/**
 * Returns a sorted shallow copy of the items array without mutating the original.
 */
export function sortItems(items: FarmItem[], sort: string): FarmItem[] {
  const sorted = [...items];
  switch (sort) {
    case "price_asc":
      return sorted.sort((a, b) => a.price - b.price);
    case "price_desc":
      return sorted.sort((a, b) => b.price - a.price);
    case "name_asc":
      return sorted.sort((a, b) => a.name.localeCompare(b.name));
    case "name_desc":
      return sorted.sort((a, b) => b.name.localeCompare(a.name));
    default:
      return sorted;
  }
}

export async function displayItems(
  type: ItemType,
  content: HTMLElement,
  isLoggedIn: boolean,
  options: DisplayItemsOptions = {}
): Promise<void> {
  const { limit = 10, offset = 0, search = "", category = "", sort = "" } = options;

  const updateOptions = (overrides: Partial<DisplayItemsOptions>) =>
    displayItems(type, content, isLoggedIn, {
      limit,
      offset,
      search,
      category,
      sort,
      ...overrides,
    });

  const refresh = () => updateOptions({});

  const container = createElement("div", { class: "protoolspage" }, [
    createElement("h2", { class: "page-title" }, [`${capitalize(type)}s`]),
  ]);

  const chipsWrapper = createElement("div", { class: "chips-wrapper" });
  container.appendChild(chipsWrapper);

  const { sortSelect, searchInput } = renderSearchAndSortUI(
    type,
    sort,
    search,
    (newSort, newSearch) => updateOptions({ sort: newSort, search: newSearch, offset: 0 })
  );

  const topBarChildren: (HTMLElement | null)[] = [
    searchInput,
    sortSelect,
    isLoggedIn
      ? Button({
          title: `Create ${type}`,
          id: `create-${type}-btn`,
          classes: "primary-button critical-action",
          events: {
            click: () => openItemFormModal("create", null, type, refresh),
          },
        })
      : null,
  ];

  container.appendChild(
    createElement("div", { class: "items-topbar" }, topBarChildren.filter(Boolean) as HTMLElement[])
  );

  // Trigger category chips rendering & items data fetching in parallel
  const chipsPromise = renderCategoryChips(
    chipsWrapper,
    category,
    (newCategory) => updateOptions({ category: newCategory, offset: 0 }),
    type
  );

  try {
    const [result] = await Promise.all([
      fetchFarmItems(type, { limit, offset, search, category }),
      chipsPromise,
    ]);

    const rawItems = result.items || [];
    const total = result.total ?? rawItems.length;

    if (rawItems.length === 0) {
      container.appendChild(createElement("p", { class: "no-results" }, [`No ${type}s found.`]));
    } else {
      const sortedItems = sortItems(rawItems, sort);
      const grid = createElement("div", { class: `${type}-grid items-grid` });

      sortedItems.forEach((item) => {
        grid.appendChild(renderItemCard(item, type, isLoggedIn, container, refresh));
      });

      container.appendChild(grid);

      renderPagination(container, total, limit, offset, (currentPage) =>
        updateOptions({ offset: (currentPage - 1) * limit })
      );
    }
  } catch (err) {
    console.error(`Failed to load ${type}s:`, err);
    container.appendChild(createElement("p", { class: "error-message" }, [`Failed to load ${type}s.`]));
  }

  content.replaceChildren(container);
}