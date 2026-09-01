import { createUserControls } from "../farm/displayFarmHelpers";
import { createElement } from "../../../components/createElement";
import { apiFetch } from "../../../api/api";
import { navigate } from "../../../routes/navigate";
import Notify from "../../../components/ui/Notify";
import Button from "../../../components/base/Button";
import { createFilterForm } from "./components/filterForm";
import { createListingCard } from "./components/listingCard";
import type { CropApiResponse, CropListing } from "./types";

/**
 * Main entry function to fetch and display crop listings.
 */
export async function displayCrop(content: HTMLElement, cropID: string | number, isLoggedIn: boolean): Promise<void> {
  const container = createElement("div", { class: "croppage" }) as HTMLElement;
  content.replaceChildren(container);

  try {
    const resp = await apiFetch<CropApiResponse>(`/crops/crop/${cropID}?page=1&limit=100`);
    if (!resp?.success || !Array.isArray(resp?.listings) || resp.listings.length === 0) {
      Notify("No listings found for this crop.", { type: "error", dismissible: true });
      return;
    }

    const listings = resp.listings as CropListing[];

    // 1. Header UI
    const header = createElement("header", { class: "crop-header" }, [
      createElement(
        "h1",
        {
          class: "crop-title",
          events: { click: () => navigate(`/aboutcrop/${cropID}`) },
          style: { fontSize: "2rem", cursor: "pointer" }
        },
        [`${resp.name || "Crop"} (${resp.category || "Uncategorized"})`]
      ),
      createElement("p", { class: "crop-meta" }, [`Total Listings: ${resp.total ?? listings.length}`])
    ]) as HTMLElement;

    // 2. Setup Filters & Listings Wrapper
    const toggleFiltersBtn = Button({ title: "Filters", id: "button", events: {}, classes: "toggle-filters-btn buttonx" }) as HTMLElement;
    const filterForm = createFilterForm();
    const listingsWrapper = createElement("section", { class: "crop-listings" }) as HTMLElement;

    // 3. Render Handler
    const renderListings = (data: CropListing[]): void => {
      listingsWrapper.replaceChildren();
      if (!data || data.length === 0) {
        listingsWrapper.appendChild(
          createElement("p", { class: "no-results" }, ["No listings match the selected filters."]) as HTMLElement
        );
        return;
      }

      const fragment = document.createDocumentFragment();
      data.forEach((listing) => {
        fragment.appendChild(createListingCard(listing, resp.name || "Crop", isLoggedIn));
      });
      listingsWrapper.appendChild(fragment);
    };

    // Initial Population
    renderListings(listings);

    // 4. Interaction Binding
    // Wire interactions
    // lazy import to keep controllers isolated
    const { setupFilterInteractions } = await import("./controllers/setupFilterInteractions");
    setupFilterInteractions({ filterForm, toggleFiltersBtn, listings, onFiltered: renderListings });

    container.append(header, toggleFiltersBtn, filterForm, listingsWrapper);
  } catch (err: unknown) {
    const errorMessage = err instanceof Error ? err.message : "Failed to load crop details.";
    Notify(errorMessage, { type: "error", dismissible: true });
  }
}
