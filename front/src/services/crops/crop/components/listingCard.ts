import { createElement } from "../../../../components/createElement";
import Imagex from "../../../../components/base/Imagex";
import { resolveImagePath, PictureType, EntityType } from "../../../../utils/imagePaths";
import { createUserControls } from "../../farm/displayFarmHelpers";
import { formatRelativeDate, formatAvailability, getStockStatus } from "../controllers/helpers";
import type { CropListing } from "../types";
import { navigate } from "../../../../routes/navigate";

export function createListingCard(listing: CropListing, cropName: string, isLoggedIn: boolean): HTMLElement {
  const imageSrc = resolveImagePath(EntityType.CROP, PictureType.THUMB, listing?.banner);
  const farmName = listing?.farmName || "Unnamed Farm";

  const imageSection = createElement("div", { class: "listing-image" }, [
    Imagex({ src: imageSrc, alt: listing?.breed || farmName, loading: "lazy" })
  ]);

  const detailRows = [
    createElement("h3", { class: "farm-link" }, [
      createElement(
        "a",
        { events: { click: () => navigate(`/farm/${listing.farmid}`) } },
        [farmName]
      )
    ]),
    createElement("p", {}, [`Breed: ${listing?.breed || "Not specified"}`]),
    createElement("p", {}, [`Location: ${listing?.location || "Unknown"}`]),
    createElement("p", {}, [`Price: ₹${Number(listing?.pricePerKg || 0).toLocaleString()}/${listing?.unit || "kg"}`]),
    createElement("p", {}, [`Available: ${listing?.availableQtyKg ?? 0} ${listing?.unit || "kg"}`]),
    createElement("p", {}, [`Inventory Value: ₹${Number(listing?.inventoryValue || 0).toLocaleString()}`]),
    createElement("p", {}, [`Status: ${listing?.outOfStock ? "Out of Stock" : getStockStatus(listing?.availableQtyKg || 0)}`]),
    createElement("p", {}, [`Featured: ${listing?.featured ? "Yes" : "No"}`]),
    createElement("p", {}, [`Rating: ${listing?.avgRating || 0} (${listing?.reviewCount || 0} reviews)`]),
    createElement("p", {}, [`Favorites: ${listing?.favoritesCount || 0}`]),
    createElement("p", {}, [`Harvest Date: ${listing?.harvestDate ? new Date(listing.harvestDate).toLocaleDateString() : "N/A"}`]),
    createElement("p", {}, [`Planted Date: ${listing?.plantedDate ? new Date(listing.plantedDate).toLocaleDateString() : "N/A"}`]),
    createElement("p", {}, [`Last Sold: ${listing?.lastSoldAt ? formatRelativeDate(listing.lastSoldAt) : "Never"}`]),
    createElement("p", {}, [`Availability: ${formatAvailability(listing?.availability)}`]),
    createElement("p", {}, [`Phone: ${listing?.phone || "N/A"}`]),
    listing?.tags?.length ? createElement("p", {}, [`Tags: ${listing.tags.join(", ")}`]) : null
  ].filter(Boolean) as HTMLElement[];

  const detailsSection = createElement("div", { class: "listing-details" }, detailRows);

  const cropData = {
    name: cropName,
    cropid: listing?.cropid,
    pricePerKg: listing?.pricePerKg,
    unit: "kg",
    breed: listing?.breed,
    quantity: listing?.availableQtyKg ?? 0
  };

  const controls = createUserControls(
    cropData,
    farmName,
    listing?.farmid,
    isLoggedIn
  );

  const controlsSection = createElement("div", { class: "listing-controls" }, controls);

  return createElement("div", { class: "listing-card" }, [
    imageSection,
    createElement("div", { class: "listing-content" }, [detailsSection, controlsSection])
  ]) as HTMLElement;
}
