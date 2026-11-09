import { createElement } from "../../../components/createElement.js";
import Modal from "../../../components/ui/Modal.js";

/* =========================================================
   TYPES & INTERFACES
   ========================================================= */

export interface MandiDetail {
    market: string;
    district: string;
    state: string;
    min: number;
    modal: number;
    max: number;
    updated: string;
}

export interface PriceItem {
    crop: string;
    category: string;
    market: string;
    district: string;
    state: string;
    min: number;
    modal: number;
    max: number;
    trend: number;
    updated: string;
    mandis: MandiDetail[];
}

interface AppState {
    search: string;
    category: string;
    view: "grid" | "list";
    sortBy: string;
}

/* =========================================================
   MOCK DATA WITH COMPREHENSIVE REGIONAL COVERAGE
   ========================================================= */

const PRICE_DATA: PriceItem[] = [
    { 
        crop: "Wheat", 
        category: "Grains", 
        market: "Pune APMC", 
        district: "Pune", 
        state: "Maharashtra", 
        min: 2350, 
        modal: 2450, 
        max: 2550, 
        trend: 2.4, 
        updated: "10 mins ago",
        mandis: [
            { market: "Pune APMC", district: "Pune", state: "Maharashtra", min: 2350, modal: 2450, max: 2550, updated: "10 mins ago" },
            { market: "Nashik Grain Market", district: "Nashik", state: "Maharashtra", min: 2300, modal: 2400, max: 2500, updated: "25 mins ago" },
            { market: "Indore Grain Mandi", district: "Indore", state: "Madhya Pradesh", min: 2280, modal: 2380, max: 2480, updated: "15 mins ago" },
            { market: "Azadpur Mandi", district: "Delhi", state: "Delhi", min: 2400, modal: 2500, max: 2620, updated: "5 mins ago" }
        ]
    },
    { 
        crop: "Onion", 
        category: "Vegetables", 
        market: "Lasalgaon", 
        district: "Nashik", 
        state: "Maharashtra", 
        min: 1650, 
        modal: 1800, 
        max: 2000, 
        trend: -5.1, 
        updated: "6 mins ago",
        mandis: [
            { market: "Lasalgaon APMC", district: "Nashik", state: "Maharashtra", min: 1650, modal: 1800, max: 2000, updated: "6 mins ago" },
            { market: "Pimpalgaon Mandi", district: "Nashik", state: "Maharashtra", min: 1600, modal: 1750, max: 1950, updated: "18 mins ago" },
            { market: "Kolkata Wholesale", district: "Kolkata", state: "West Bengal", min: 1900, modal: 2100, max: 2300, updated: "30 mins ago" },
            { market: "Bangalore City Market", district: "Bangalore", state: "Karnataka", min: 1850, modal: 2000, max: 2200, updated: "12 mins ago" }
        ]
    },
    { 
        crop: "Tomato", 
        category: "Vegetables", 
        market: "Delhi Mandi", 
        district: "Delhi", 
        state: "Delhi", 
        min: 900, 
        modal: 1200, 
        max: 1500, 
        trend: 1.6, 
        updated: "12 mins ago",
        mandis: [
            { market: "Azadpur Mandi", district: "Delhi", state: "Delhi", min: 900, modal: 1200, max: 1500, updated: "12 mins ago" },
            { market: "Kolar Vegetable Market", district: "Kolar", state: "Karnataka", min: 800, modal: 1050, max: 1300, updated: "22 mins ago" },
            { market: "Vashi APMC", district: "Navi Mumbai", state: "Maharashtra", min: 950, modal: 1250, max: 1550, updated: "8 mins ago" }
        ]
    },
    { 
        crop: "Cotton", 
        category: "Oilseeds", 
        market: "Akola", 
        district: "Akola", 
        state: "Maharashtra", 
        min: 7000, 
        modal: 7200, 
        max: 7400, 
        trend: 0, 
        updated: "20 mins ago",
        mandis: [
            { market: "Akola APMC", district: "Akola", state: "Maharashtra", min: 7000, modal: 7200, max: 7400, updated: "20 mins ago" },
            { market: "Rajkot Ginned Market", district: "Rajkot", state: "Gujarat", min: 7100, modal: 7300, max: 7500, updated: "40 mins ago" }
        ]
    },
    { 
        crop: "Soybean", 
        category: "Legumes", 
        market: "Indore", 
        district: "Indore", 
        state: "Madhya Pradesh", 
        min: 4800, 
        modal: 4950, 
        max: 5100, 
        trend: 3.2, 
        updated: "8 mins ago",
        mandis: [
            { market: "Indore Mandi", district: "Indore", state: "Madhya Pradesh", min: 4800, modal: 4950, max: 5100, updated: "8 mins ago" },
            { market: "Ujjain Grain Market", district: "Ujjain", state: "Madhya Pradesh", min: 4750, modal: 4900, max: 5050, updated: "15 mins ago" },
            { market: "Kota Mandi", district: "Kota", state: "Rajasthan", min: 4850, modal: 5000, max: 5150, updated: "25 mins ago" }
        ]
    }
];

const state: AppState = {
    search: "",
    category: "All",
    view: "grid", 
    sortBy: "default" 
};

let appContainer: HTMLElement | null = null;
let gridContainer: HTMLElement | null = null;
let searchInputRef: HTMLInputElement | null = null;

const categories: string[] = ["All", ...Array.from(new Set(PRICE_DATA.map(x => x.category)))];
const trendingCrops: string[] = ["Onion", "Wheat", "Tomato"];

/* =========================================================
   ADVANCED FILTER & REGIONAL SEARCH LOGIC
   ========================================================= */

function getProcessedData(): PriceItem[] {
    const query = state.search.trim().toLowerCase();
    
    const result = PRICE_DATA.filter(item => {
        // Matches main crop name, or category
        const matchesMain = item.crop.toLowerCase().includes(query) || 
                            item.category.toLowerCase().includes(query);
        
        // Matches across any nested mandi state, district, or market name
        const matchesMandiRegion = item.mandis.some(m => 
            m.state.toLowerCase().includes(query) ||
            m.district.toLowerCase().includes(query) ||
            m.market.toLowerCase().includes(query)
        );

        const matchesSearch = !query || matchesMain || matchesMandiRegion;
        const matchesCategory = state.category === "All" || item.category === state.category;
        
        return matchesSearch && matchesCategory;
    });

    if (state.sortBy === "high-low") {
        result.sort((a, b) => b.modal - a.modal);
    } else if (state.sortBy === "low-high") {
        result.sort((a, b) => a.modal - b.modal);
    } else if (state.sortBy === "trend") {
        result.sort((a, b) => Math.abs(b.trend) - Math.abs(a.trend));
    }

    return result;
}

/* =========================================================
   MODAL COMPONENT BUILDER
   ========================================================= */

function openMandiModal(item: PriceItem): void {
    const modalSubtitle = createElement("p", { class: "mp-modal-subtitle" }, [
        `Category: ${item.category} • Showing live wholesale rates by state & district`
    ]);

    const tableHeader = createElement("div", { class: "mp-modal-tr mp-modal-th" }, [
        createElement("span", {}, ["Market / Mandi"]),
        createElement("span", {}, ["District, State"]),
        createElement("span", {}, ["Min (₹)"]),
        createElement("span", {}, ["Modal Rate (₹)"]),
        createElement("span", {}, ["Max (₹)"]),
        createElement("span", {}, ["Updated"])
    ]);

    const tableRows = item.mandis.map(m => createElement("div", { class: "mp-modal-tr" }, [
        createElement("strong", {}, [m.market]),
        createElement("span", {}, [`${m.district}, ${m.state}`]),
        createElement("span", {}, [`₹${m.min.toLocaleString("en-IN")}`]),
        createElement("span", { class: "mp-modal-highlight-rate" }, [`₹${m.modal.toLocaleString("en-IN")}`]),
        createElement("span", {}, [`₹${m.max.toLocaleString("en-IN")}`]),
        createElement("span", { class: "mp-time" }, [m.updated])
    ]));

    const tableContainer = createElement("div", { class: "mp-modal-table" }, [tableHeader, ...tableRows]);
    const modalContentBox = createElement("div", { class: "mp-modal-content" }, [modalSubtitle, tableContainer]);

    Modal({
        title: `${item.crop} Prices Across Indian Mandis`,
        content: modalContentBox,
        size: "large",
        closeOnOverlayClick: true,
        showHeader: true,
        showCloseButton: true,
        autofocus: false
    });
}

/* =========================================================
   UI COMPONENT BUILDERS
   ========================================================= */

function createHeaderComponent(): HTMLElement {
    const title = createElement("h1", { class: "mp-title" }, ["Live Mandi Market Prices"]);
    
    searchInputRef = createElement("input", {
        type: "text",
        placeholder: "Search crops, states (e.g. Maharashtra), districts or markets...",
        class: "mp-search-input",
        value: state.search,
        events: {
            input: (e: Event) => {
                const target = e.target as HTMLInputElement;
                state.search = target.value;
                renderGrid();
            }
        }
    }) as HTMLInputElement;

    const suggestionLabel = createElement("span", { class: "mp-suggest-label" }, ["Trending:"]);
    const suggestionContainer = createElement("div", { class: "mp-suggest-box" }, [
        suggestionLabel,
        ...trendingCrops.map(crop => createElement("button", {
            type: "button",
            class: "mp-suggest-pill",
            events: {
                click: () => {
                    state.search = crop;
                    if (searchInputRef) searchInputRef.value = crop;
                    renderGrid();
                }
            }
        }, [crop]))
    ]);

    const searchWrapper = createElement("div", { class: "mp-search-wrapper" }, [searchInputRef, suggestionContainer]);

    const sortDropdown = createElement("select", {
        class: "mp-sort-select",
        ariaLabel: "Sort Mandi Prices",
        events: {
            change: (e: Event) => {
                const target = e.target as HTMLSelectElement;
                state.sortBy = target.value;
                renderGrid();
            }
        }
    }, [
        createElement("option", { value: "default" }, ["Sort: Default"]),
        createElement("option", { value: "high-low" }, ["Price: High to Low"]),
        createElement("option", { value: "low-high" }, ["Price: Low to High"]),
        createElement("option", { value: "trend" }, ["Highest Volatility"])
    ]);
    (sortDropdown as HTMLSelectElement).value = state.sortBy;

    const tabsWrapper = createElement("div", { class: "mp-tabs-wrapper" });
    categories.forEach(cat => {
        const btn = createElement("button", {
            type: "button",
            class: `mp-filter-btn ${state.category === cat ? "active" : ""}`,
            events: {
                click: () => {
                    state.category = cat;
                    document.querySelectorAll(".mp-filter-btn").forEach(b => b.classList.remove("active"));
                    btn.classList.add("active");
                    renderGrid();
                }
            }
        }, [cat]);
        tabsWrapper.appendChild(btn);
    });

    const toggleBtn = createElement("button", {
        type: "button",
        class: "mp-view-toggle",
        events: {
            click: () => {
                state.view = state.view === "grid" ? "list" : "grid";
                toggleBtn.textContent = state.view === "grid" ? "Horizontal List View" : "Compact Grid View";
                renderGrid();
            }
        }
    }, [state.view === "grid" ? "Horizontal List View" : "Compact Grid View"]);
    
    const controlBar = createElement("div", { class: "mp-control-bar" }, [tabsWrapper, sortDropdown, toggleBtn]);

    return createElement("div", { class: "mp-dashboard-header" }, [title, searchWrapper, controlBar]);
}

function createPriceCard(item: PriceItem): HTMLElement {
    const isPositive = item.trend >= 0;
    const trendClass = item.trend === 0 ? "neutral" : isPositive ? "positive" : "negative";
    const trendIcon = item.trend === 0 ? "•" : isPositive ? "▲" : "▼";
    const trendPrefix = isPositive && item.trend > 0 ? "+" : "";

    const badge = createElement("span", { class: "mp-badge" }, [item.category]);
    const time = createElement("span", { class: "mp-time" }, [item.updated]);
    const cardHeader = createElement("div", { class: "mp-card-meta" }, [badge, time]);

    const cropName = createElement("h3", { class: "mp-crop-name" }, [item.crop]);
    const location = createElement("p", { class: "mp-location-text" }, [`📍 ${item.market}, ${item.district} (${item.state})`]);
    const infoSec = createElement("div", { class: "mp-info-section" }, [cropName, location]);

    const label = createElement("span", { class: "mp-label" }, ["Avg Market Rate"]);
    const modalPrice = createElement("div", { class: "mp-modal-price" }, [`₹${item.modal.toLocaleString("en-IN")}`]);
    const trendPill = createElement("span", { class: `mp-trend-pill ${trendClass}` }, [`${trendIcon} ${trendPrefix}${item.trend}%`]);
    const primaryPriceWrapper = createElement("div", { class: "mp-modal-wrapper" }, [label, modalPrice, trendPill]);

    const totalRange = item.max - item.min;
    const directOffset = item.modal - item.min;
    const trackingPercentage = totalRange > 0 ? Math.min(Math.max((directOffset / totalRange) * 100, 0), 100) : 50;

    const visualTrackIndicator = createElement("div", { 
        class: "mp-range-fill-head", 
        style: `left: ${trackingPercentage}%;` 
    });
    const customProgressBarRangeLine = createElement("div", { class: "mp-range-track-line" }, [visualTrackIndicator]);
    
    const minBound = createElement("span", {}, [`₹${item.min.toLocaleString("en-IN")}`]);
    const maxBound = createElement("span", {}, [`₹${item.max.toLocaleString("en-IN")}`]);
    const progressLabelContainer = createElement("div", { class: "mp-range-labels-wrapper" }, [minBound, maxBound]);
    
    const contextRangeGaugeComponent = createElement("div", { class: "mp-range-slider-component" }, [
        createElement("span", { class: "mp-label" }, ["Market Spread Location"]),
        customProgressBarRangeLine,
        progressLabelContainer
    ]);

    const priceMatrix = createElement("div", { class: "mp-price-matrix" }, [primaryPriceWrapper, contextRangeGaugeComponent]);

    const card = createElement("div", { 
        class: `mp-price-card ${state.view}-layout`,
        title: "Click to view prices across all Mandis",
        events: {
            click: () => openMandiModal(item)
        }
    }, [cardHeader, infoSec, priceMatrix]);

    return card;
}

function createEmptyState(): HTMLElement {
    return createElement("button", { 
        type: "button",
        class: "mp-empty-state-reset",
        events: {
            click: () => {
                state.search = "";
                state.category = "All";
                state.sortBy = "default";
                if (searchInputRef) searchInputRef.value = "";
                const firstFilterBtn = document.querySelector(".mp-filter-btn");
                document.querySelectorAll(".mp-filter-btn").forEach(b => b.classList.remove("active"));
                if (firstFilterBtn) firstFilterBtn.classList.add("active");
                renderGrid();
            }
        }
    }, ["No crops or states matched your search query. Click here to reset filters."]);
}

function renderGrid(): void {
    if (!gridContainer) return;
    
    gridContainer.className = `mp-grid-container layout-${state.view}`;
    gridContainer.innerHTML = "";

    const activeDataset = getProcessedData();

    if (activeDataset.length === 0) {
        gridContainer.appendChild(createEmptyState());
        return;
    }

    activeDataset.forEach(item => {
        gridContainer?.appendChild(createPriceCard(item));
    });
}

export function displayBazarBhav(rootContainer: HTMLElement | null): void {
    if (!rootContainer) return;
    appContainer = rootContainer;
    appContainer.innerHTML = ""; 
    
    const dashboardHeader = createHeaderComponent();
    gridContainer = createElement("div", { class: `mp-grid-container layout-${state.view}` });

    appContainer.appendChild(dashboardHeader);
    appContainer.appendChild(gridContainer);

    renderGrid();
}