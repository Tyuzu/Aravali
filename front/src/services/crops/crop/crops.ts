import { createElement } from "../../../components/createElement.js";
import { apiFetch } from "../../../api/api.js";
import { guessCategoryFromName } from "./displayCropshelpers.js";
import { navigate } from "../../../routes/navigate.js";
import { renderCropCard } from "./components/cropCard.js";
import { debounce } from "../../../utils/deutils.js";
import Button from "../../../components/base/Button.js";
import { createMainLayout } from "../../../components/layout/mainLayout.js";
import { createAsideContent } from "../../../components/layout/asideLayout.js";

// --- Types & Interfaces ---

export interface Crop {
  name: string;
  minPrice?: number;
  maxPrice?: number;
  availableCount?: number;
  unit?: string;
  banner?: string;
  tags?: string[];
  seasonMonths?: number[];
  price?: number;
  quantity?: number;
  farmName?: string;
}

export type CategorizedCrops = Record<string, Crop[]>;
export type ViewMode = "grid" | "list";

export interface MarketHighlightItem {
  id: string;
  text: string;
}

export interface MarketHighlightsData {
  activeDeals?: MarketHighlightItem[];
  seasonalPicks?: MarketHighlightItem[];
  cropTrends?: MarketHighlightItem[];
  announcements?: MarketHighlightItem[];
  showcase?: MarketHighlightItem[];
}

interface FilterOptions {
  term: string;
  tags: Set<string>;
  sortBy: string;
}

interface InterfaceState {
  cropData: CategorizedCrops;
  categories: string[];
  currentTab: string | null;
  activeTags: Set<string>;
  viewMode: ViewMode;
  searchBox: HTMLInputElement;
  sortSelect: HTMLSelectElement;
  gridViewBtn: HTMLButtonElement;
  listViewBtn: HTMLButtonElement;
  tabs: Record<string, HTMLElement>;
  tabButtons: HTMLElement;
}

interface RawCropType {
  Name?: string;
  MinPrice?: number;
  MaxPrice?: number;
  AvailableCount?: number;
  Unit?: string;
  Banner?: string;
}

interface CropsApiResponse {
  cropTypes?: RawCropType[];
}

/**
 * Creates formatted promo items/list configuration for createAsideContent sections.
 */
function createPromoSection(title: string, items: MarketHighlightItem[] = []) {
  return {
    title,
    className: "promo-box",
    content: createElement(
      "ul",
      { class: "promo-list" },
      items.length > 0 
        ? items.map((item) => createElement("li", { key: item.id }, [item.text]))
        : [createElement("li", { class: "no-data" }, ["No updates available"])]
    )
  };
}

export function cropAside(_cropData: CategorizedCrops, highlights: MarketHighlightsData = {}) {
  return createAsideContent({
    title: "Market Highlights",
    actions: [
      Button({
        title: "Buy Products",
        id: "buyprds-crp-btn",
        events: { click: () => navigate("/products") },
        classes: "action-btn buttonx primary"
      }),
      Button({
        title: "See Recipes",
        id: "recipes-crp-btn",
        events: { click: () => navigate("/recipes") },
        classes: "buttonx secondary"
      }),
      Button({
        title: "List Your Farm",
        id: "newfrm-btn",
        events: { click: () => navigate("/create-farm") },
        classes: "buttonx secondary"
      })
    ],
    sections: [
      createPromoSection("💸 Active Deals", highlights.activeDeals),
      createPromoSection("📅 Seasonal Picks", highlights.seasonalPicks),
      createPromoSection("📊 Crop Trends", highlights.cropTrends),
      createPromoSection("🔔 Announcements", highlights.announcements),
      createPromoSection("📷 Farmer's Showcase", highlights.showcase)
    ],
    showAd: true,
    adOptions: {
      layout: "vertical"
    }
  });
}

// --- Helpers & Utils ---

function filterAndSortCrops(crops: Crop[] = [], { term, tags, sortBy }: FilterOptions): Crop[] {
  const searchTerm = term.toLowerCase();

  return crops
    .filter(crop => {
      const matchesTerm = crop.name?.toLowerCase().includes(searchTerm);
      const matchesTags = [...tags].every(tag => crop.tags?.includes(tag));
      return matchesTerm && matchesTags;
    })
    .sort((a, b) =>
      sortBy === "az"
        ? a.name.localeCompare(b.name)
        : b.name.localeCompare(a.name)
    );
}

function formatCropSlug(name: string): string {
  return name.toLowerCase().replace(/\s+/g, "_");
}

// --- Interface State Management ---

export function renderCropInterface(
  container: HTMLElement, 
  cropData: CategorizedCrops, 
  highlights: MarketHighlightsData
): void {
  const mainContent = createElement("div", { class: "catalogue-main" });

  const searchBox = createElement("input", {
    type: "text",
    name: "search",
    placeholder: "Search crops by name...",
    class: "search-box",
    "aria-label": "Search crops"
  }) as HTMLInputElement;

  const sortSelect = createElement("select", { class: "sort-box", name: "sortby", "aria-label": "Sort crops" }, [
    createElement("option", { value: "az" }, ["A → Z"]),
    createElement("option", { value: "za" }, ["Z → A"])
  ]) as HTMLSelectElement;

  const gridViewBtn = createElement(
    "button",
    { class: "view-btn active", type: "button", "aria-label": "Grid View", title: "Grid View" },
    ["田"]
  ) as HTMLButtonElement;

  const listViewBtn = createElement(
    "button",
    { class: "view-btn", type: "button", "aria-label": "List View", title: "List View" },
    ["☰"]
  ) as HTMLButtonElement;

  const viewToggleGroup = createElement("div", { class: "view-toggle-group", role: "group", "aria-label": "View mode" }, [
    gridViewBtn,
    listViewBtn
  ]);

  const controls = createElement("div", { class: "top-controls" }, [
    searchBox,
    sortSelect,
    viewToggleGroup
  ]);

  const tabButtons = createElement("div", { class: "tabs", role: "tablist" });
  const tabsWrapper = createElement("div", { id: "catalogue-container" });

  mainContent.append(
    createElement("h3", {}, ["All Crops"]),
    controls,
    tabButtons,
    tabsWrapper
  );

  const categories = Object.keys(cropData);
  const state: InterfaceState = {
    cropData,
    categories,
    currentTab: categories[0] || null,
    activeTags: new Set(),
    viewMode: "grid",
    searchBox,
    sortSelect,
    gridViewBtn,
    listViewBtn,
    tabs: {},
    tabButtons
  };

  const setViewMode = (mode: ViewMode) => {
    if (state.viewMode === mode) return;
    state.viewMode = mode;

    gridViewBtn.classList.toggle("active", mode === "grid");
    listViewBtn.classList.toggle("active", mode === "list");
    updateAllTabs(state);
  };

  gridViewBtn.onclick = () => setViewMode("grid");
  listViewBtn.onclick = () => setViewMode("list");

  categories.forEach((cat, index) => {
    const isFirst = index === 0;
    const count = cropData[cat]?.length || 0;
    const tabId = `tab-${formatCropSlug(cat)}`;

    const btn = createElement(
      "button",
      {
        class: `buttonx ${isFirst ? "active" : ""}`,
        role: "tab",
        id: tabId,
        "aria-selected": String(isFirst),
        "aria-controls": `panel-${formatCropSlug(cat)}`,
        disabled: count === 0,
        "data-category": cat
      },
      [`${cat.charAt(0).toUpperCase() + cat.slice(1)} (${count})`]
    ) as HTMLButtonElement;

    btn.onclick = () => {
      state.currentTab = cat;
      updateAllTabs(state);
    };

    tabButtons.appendChild(btn);

    const pane = createElement("div", {
      class: `tab-content ${state.viewMode}-view`,
      role: "tabpanel",
      id: `panel-${formatCropSlug(cat)}`,
      "aria-labelledby": tabId
    });
    state.tabs[cat] = pane;
    tabsWrapper.appendChild(pane);
  });

  sortSelect.onchange = () => updateAllTabs(state);
  searchBox.addEventListener("input", debounce(() => updateAllTabs(state), 250));

  updateAllTabs(state);

  const layout = createMainLayout({
    mainContent: [mainContent],
    asideContent: cropAside(cropData, highlights),
    pageClass: "catalogue-layout",
    showMainAd: true,
    mainAdPlacement: "top"
  });

  container.appendChild(layout);
}

function updateAllTabs(state: InterfaceState): void {
  const { categories, currentTab, tabButtons, tabs, viewMode } = state;
  if (!currentTab) return;

  categories.forEach(cat => {
    const pane = tabs[cat];
    if (pane) {
      const isCurrent = cat === currentTab;
      pane.style.display = isCurrent ? "flex" : "none";
      pane.classList.remove("grid-view", "list-view");
      pane.classList.add(`${viewMode}-view`);

      if (isCurrent) {
        updateTab(cat, state);
      }
    }
  });

  Array.from(tabButtons.children).forEach(btn => {
    const htmlBtn = btn as HTMLElement;
    const btnCategory = htmlBtn.dataset.category || "";
    const isActive = btnCategory.toLowerCase() === currentTab.toLowerCase();
    
    htmlBtn.classList.toggle("active", isActive);
    htmlBtn.setAttribute("aria-selected", String(isActive));
  });
}

function updateTab(category: string, state: InterfaceState): void {
  const { cropData, tabs, searchBox, sortSelect, activeTags, viewMode } = state;
  const container = tabs[category];

  if (!container) return;
  container.replaceChildren();

  const filtered = filterAndSortCrops(cropData[category], {
    term: searchBox.value.trim(),
    tags: activeTags,
    sortBy: sortSelect.value
  });

  if (filtered.length === 0) {
    container.appendChild(
      createElement("p", { class: "empty-category" }, ["No crops found matching your criteria."])
    );
    return;
  }

  const fragment = document.createDocumentFragment();
  filtered.forEach(crop => {
    fragment.appendChild(renderCropCard(crop, viewMode));
  });
  container.appendChild(fragment);
}

// --- Main Entrypoint ---

export async function displayCrops(content: HTMLElement): Promise<void> {
  const contentContainer = createElement("div", { class: "cropspage" });
  content.replaceChildren(contentContainer);

  const categorized: CategorizedCrops = {};
  let highlights: MarketHighlightsData = {};

  try {
    // Fetch crops and market highlights concurrently
    const [cropsResponse, highlightsResponse] = await Promise.all([
      apiFetch<CropsApiResponse>("/crops/types"),
      apiFetch<MarketHighlightsData>("/market/highlights").catch(() => ({})) // Graceful fallback if endpoint fails
    ]);

    if (!cropsResponse?.cropTypes || !Array.isArray(cropsResponse.cropTypes)) {
      throw new Error("Invalid response format: 'cropTypes' array missing");
    }

    highlights = highlightsResponse || {};

    cropsResponse.cropTypes.forEach((raw: RawCropType) => {
      if (!raw.Name) return;

      const crop: Crop = {
        name: raw.Name,
        minPrice: raw.MinPrice,
        maxPrice: raw.MaxPrice,
        availableCount: raw.AvailableCount,
        unit: raw.Unit,
        banner: raw.Banner || "placeholder.jpg",
        tags: [],
        seasonMonths: []
      };

      const category = guessCategoryFromName(crop.name);
      if (!categorized[category]) {
        categorized[category] = [];
      }
      categorized[category].push(crop);
    });
  } catch (err) {
    console.error("Error fetching market data:", err);
    contentContainer.appendChild(
      createElement("div", { class: "error-banner" }, [
        createElement("p", {}, ["Failed to load marketplace data. Please try again later."])
      ])
    );
    return;
  }

  renderCropInterface(contentContainer, categorized, highlights);
}