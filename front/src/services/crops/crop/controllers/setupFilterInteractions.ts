import Notify from "../../../../components/ui/Notify";
import type { SetupFilterInteractionsParams } from "../types";
import { debounce } from "./utils";

export function setupFilterInteractions({ filterForm, toggleFiltersBtn, listings, onFiltered }: SetupFilterInteractionsParams): void {
  const inputs = {
    location: filterForm.querySelector<HTMLInputElement>("#filter-location"),
    breed: filterForm.querySelector<HTMLInputElement>("#filter-breed"),
    minPrice: filterForm.querySelector<HTMLInputElement>("#filter-min-price"),
    maxPrice: filterForm.querySelector<HTMLInputElement>("#filter-max-price"),
    minQty: filterForm.querySelector<HTMLInputElement>("#filter-min-qty"),
    maxQty: filterForm.querySelector<HTMLInputElement>("#filter-max-qty"),
    harvestDate: filterForm.querySelector<HTMLInputElement>("#filter-harvest")
  };

  const applyButton = filterForm.querySelector<HTMLButtonElement>("#apply-filters");
  const resetButton = filterForm.querySelector<HTMLButtonElement>("#reset-filters");

  if (!inputs.location || !inputs.breed || !inputs.minPrice || !inputs.maxPrice ||
    !inputs.minQty || !inputs.maxQty || !inputs.harvestDate || !applyButton || !resetButton) {
    Notify("Unable to initialize crop filters.", { type: "error", dismissible: true });
    return;
  }

  const validInputs = inputs as Record<keyof typeof inputs, HTMLInputElement>;

  const applyFilters = (): void => {
    const filters = {
      location: validInputs.location.value.trim().toLowerCase(),
      breed: validInputs.breed.value.trim().toLowerCase(),
      minPrice: parseFloat(validInputs.minPrice.value) || null,
      maxPrice: parseFloat(validInputs.maxPrice.value) || null,
      minQty: parseFloat(validInputs.minQty.value) || null,
      maxQty: parseFloat(validInputs.maxQty.value) || null,
      harvestDate: validInputs.harvestDate.value || null
    };

    if (filters.minPrice && filters.maxPrice && filters.minPrice > filters.maxPrice) {
      Notify("Invalid price range (min > max).", { type: "warning", dismissible: true });
      return;
    }
    if (filters.minQty && filters.maxQty && filters.minQty > filters.maxQty) {
      Notify("Invalid quantity range (min > max).", { type: "warning", dismissible: true });
      return;
    }

    const filteredListings = listings.filter((listing) => {
      const locationMatch = !filters.location || (listing?.location || "").toLowerCase().includes(filters.location as string);
      const breedMatch = !filters.breed || (listing?.breed || "").toLowerCase().includes(filters.breed as string);

      const priceMatch =
        (!filters.minPrice || (listing?.pricePerKg ?? 0) >= filters.minPrice) &&
        (!filters.maxPrice || (listing?.pricePerKg ?? 0) <= filters.maxPrice);

      const qtyMatch =
        (!filters.minQty || (listing?.availableQtyKg ?? 0) >= filters.minQty) &&
        (!filters.maxQty || (listing?.availableQtyKg ?? 0) <= filters.maxQty);

      let harvestMatch = true;
      if (filters.harvestDate) {
        if (!listing?.harvestDate) {
          harvestMatch = false;
        } else {
          const parsed = new Date(listing.harvestDate);
          harvestMatch = !isNaN(parsed.getTime()) && parsed.toISOString().split("T")[0] === filters.harvestDate;
        }
      }

      return locationMatch && breedMatch && priceMatch && qtyMatch && harvestMatch;
    });

    onFiltered(filteredListings);
  };

  const debouncedApply = debounce(applyFilters, 250);

  Object.values(validInputs).forEach((input) => {
    input.addEventListener("input", debouncedApply);
  });

  const resetFilters = (): void => {
    filterForm.reset();
    onFiltered(listings);
    filterForm.classList.remove("open");
  };

  toggleFiltersBtn.addEventListener("click", () => filterForm.classList.toggle("open"));
  applyButton.addEventListener("click", () => {
    applyFilters();
    filterForm.classList.remove("open");
  });
  resetButton.addEventListener("click", resetFilters);

  filterForm.addEventListener("keydown", (e: KeyboardEvent) => {
    if (e.key === "Enter") {
      e.preventDefault();
      applyFilters();
      filterForm.classList.remove("open");
    }
  });
}
